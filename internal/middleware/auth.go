package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"google-indexing-api/internal/config"
	"google-indexing-api/internal/models"
)

func APIKeyAuth(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()

		// Skip auth for health check
		if c.Request.URL.Path == "/api/health" {
			c.Next()
			return
		}

		// Get API key from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Missing Authorization header",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		// Extract API key from Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid Authorization header format",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		apiKey := parts[1]
		if apiKey != cfg.API.Key {
			logger.WithField("provided_key", apiKey).Warn("Invalid API key")
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "Unauthorized",
				Message: "Invalid API key",
				Code:    http.StatusUnauthorized,
			})
			c.Abort()
			return
		}

		logger.Debug("API key authentication successful")
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()

		origin := c.Request.Header.Get("Origin")

		// Simple CORS - in production, implement proper origin checking
		if len(cfg.CORS.AllowedOrigins) > 0 && cfg.CORS.AllowedOrigins[0] == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"status":     param.StatusCode,
			"method":     param.Method,
			"path":       param.Path,
			"ip":         param.ClientIP,
			"user_agent": param.Request.UserAgent(),
			"latency":    param.Latency,
		}).Info("HTTP Request")

		return ""
	})
}

func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.WithError(err).Error("Request processing error")

			if !c.Writer.Written() {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error:   "Internal Server Error",
					Message: err.Error(),
					Code:    http.StatusInternalServerError,
				})
			}
		}
	}
}
