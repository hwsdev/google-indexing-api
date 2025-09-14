package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"google-indexing-api/internal/config"
	"google-indexing-api/internal/models"
	"google-indexing-api/internal/services"
)

type IndexingHandler struct {
	service   *services.GoogleIndexingService
	logger    *logrus.Logger
	validator *validator.Validate
}

func NewIndexingHandler(service *services.GoogleIndexingService, logger *logrus.Logger) *IndexingHandler {
	return &IndexingHandler{
		service:   service,
		logger:    logger,
		validator: validator.New(),
	}
}

// @Summary Submit a single URL for indexing
// @Description Submit a single URL to Google Indexing API with service account credentials
// @Tags indexing
// @Accept json
// @Produce json
// @Param request body models.IndexRequest true "URL to index with service account"
// @Success 200 {object} models.IndexResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/index [post]
func (h *IndexingHandler) SubmitURL(c *gin.Context) {
	var req models.IndexRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate URL format
	if !h.isValidURL(req.URL) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid URL format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate service account (now required)
	if req.ServiceAccount == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Service account is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.validateServiceAccount(req.ServiceAccount); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: fmt.Sprintf("Invalid service account: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.service.SubmitURL(c.Request.Context(), req.URL, req.ServiceAccount)
	if err != nil {
		h.logger.WithError(err).Error("Failed to submit URL")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to submit URL to Google Indexing API",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusBadRequest, response)
	}
}

// @Summary Submit multiple URLs for indexing
// @Description Submit multiple URLs to Google Indexing API in batch with service account credentials
// @Tags indexing
// @Accept json
// @Produce json
// @Param request body models.BatchIndexRequest true "URLs to index with service account"
// @Success 200 {object} models.BatchIndexResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/index/batch [post]
func (h *IndexingHandler) SubmitURLsBatch(c *gin.Context) {
	var req models.BatchIndexRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.WithError(err).Error("Failed to bind JSON request")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid request format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate all URLs
	for _, urlStr := range req.URLs {
		if !h.isValidURL(urlStr) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "Bad Request",
				Message: "One or more URLs have invalid format",
				Code:    http.StatusBadRequest,
			})
			return
		}
	}

	// Limit batch size
	cfg := config.GetConfig()
	if len(req.URLs) > cfg.Performance.MaxBatchSize {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: fmt.Sprintf("Batch size cannot exceed %d URLs", cfg.Performance.MaxBatchSize),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Validate service account (now required)
	if req.ServiceAccount == nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Service account is required",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if err := h.validateServiceAccount(req.ServiceAccount); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: fmt.Sprintf("Invalid service account: %v", err),
			Code:    http.StatusBadRequest,
		})
		return
	}

	response, err := h.service.SubmitURLsBatch(c.Request.Context(), req.URLs, req.ServiceAccount)
	if err != nil {
		h.logger.WithError(err).Error("Failed to submit batch URLs")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to submit URLs to Google Indexing API",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get URL indexing status
// @Description Get the indexing status of a URL from Google
// @Tags indexing
// @Produce json
// @Param url path string true "URL to check status"
// @Success 200 {object} models.StatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/status/{url} [get]
func (h *IndexingHandler) GetURLStatus(c *gin.Context) {
	urlParam := c.Param("url")

	// URL decode the parameter
	decodedURL, err := url.QueryUnescape(urlParam)
	if err != nil {
		h.logger.WithError(err).Error("Failed to decode URL parameter")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid URL parameter",
			Code:    http.StatusBadRequest,
		})
		return
	}

	if !h.isValidURL(decodedURL) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid URL format",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// For status check, we'll use default service account
	// In future, this could be enhanced to accept service account in query params
	response, err := h.service.GetURLStatus(c.Request.Context(), decodedURL, nil)
	if err != nil {
		h.logger.WithError(err).Error("Failed to get URL status")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to get URL status from Google Indexing API",
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} models.HealthResponse
// @Router /api/health [get]
func (h *IndexingHandler) HealthCheck(c *gin.Context) {
	cfg := config.GetConfig()
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   cfg.App.Version,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get service cache statistics
// @Description Get statistics about cached service accounts
// @Tags indexing
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cache/stats [get]
func (h *IndexingHandler) GetCacheStats(c *gin.Context) {
	stats := h.service.GetCacheStats()
	c.JSON(http.StatusOK, stats)
}

// @Summary Clear service cache
// @Description Clear all cached service accounts
// @Tags indexing
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/cache/clear [post]
func (h *IndexingHandler) ClearCache(c *gin.Context) {
	h.service.ClearCache()
	c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Cache cleared successfully",
	})
}

func (h *IndexingHandler) isValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != "" && (u.Scheme == "http" || u.Scheme == "https")
}

func (h *IndexingHandler) validateServiceAccount(sa *models.ServiceAccountCredentials) error {
	if sa == nil {
		return fmt.Errorf("service account cannot be nil")
	}

	if sa.Type != "service_account" {
		return fmt.Errorf("invalid service account type, expected 'service_account'")
	}

	if sa.ProjectID == "" {
		return fmt.Errorf("project_id is required")
	}

	if sa.PrivateKey == "" {
		return fmt.Errorf("private_key is required")
	}

	if sa.ClientEmail == "" {
		return fmt.Errorf("client_email is required")
	}

	if sa.ClientID == "" {
		return fmt.Errorf("client_id is required")
	}

	if sa.AuthURI == "" {
		return fmt.Errorf("auth_uri is required")
	}

	if sa.TokenURI == "" {
		return fmt.Errorf("token_uri is required")
	}

	// Validate private key format
	if !h.isValidPrivateKey(sa.PrivateKey) {
		return fmt.Errorf("invalid private key format")
	}

	return nil
}

func (h *IndexingHandler) isValidPrivateKey(privateKey string) bool {
	// Basic validation for PEM format
	return len(privateKey) > 50 &&
		(strings.Contains(privateKey, "-----BEGIN PRIVATE KEY-----") ||
			strings.Contains(privateKey, "-----BEGIN RSA PRIVATE KEY-----"))
}
