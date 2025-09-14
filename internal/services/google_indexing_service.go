package services

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/indexing/v3"
	"google.golang.org/api/option"

	"google-indexing-api/internal/models"
)

type GoogleIndexingService struct {
	defaultService *indexing.Service
	logger         *logrus.Logger
	serviceCache   map[string]*indexing.Service
	cacheMutex     sync.RWMutex
}

func NewGoogleIndexingService(logger *logrus.Logger) (*GoogleIndexingService, error) {
	return &GoogleIndexingService{
		defaultService: nil, // No default service account
		logger:         logger,
		serviceCache:   make(map[string]*indexing.Service),
		cacheMutex:     sync.RWMutex{},
	}, nil
}

func (gis *GoogleIndexingService) getIndexingService(ctx context.Context, serviceAccount *models.ServiceAccountCredentials) (*indexing.Service, error) {
	// Service account is now required
	if serviceAccount == nil {
		return nil, fmt.Errorf("service account is required")
	}

	// Create cache key from service account
	cacheKey := serviceAccount.ClientEmail

	// Check cache first
	gis.cacheMutex.RLock()
	if cachedService, exists := gis.serviceCache[cacheKey]; exists {
		gis.cacheMutex.RUnlock()
		return cachedService, nil
	}
	gis.cacheMutex.RUnlock()

	// Create new service from provided credentials
	credentialsJSON, err := json.Marshal(serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service account credentials: %v", err)
	}

	service, err := indexing.NewService(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create indexing service with provided credentials: %v", err)
	}

	// Cache the service
	gis.cacheMutex.Lock()
	gis.serviceCache[cacheKey] = service
	gis.cacheMutex.Unlock()

	gis.logger.WithField("service_account", serviceAccount.ClientEmail).Info("Created new indexing service")

	return service, nil
}

func (gis *GoogleIndexingService) SubmitURL(ctx context.Context, url string, serviceAccount *models.ServiceAccountCredentials) (*models.IndexResponse, error) {
	gis.logger.WithField("url", url).Info("Submitting URL to Google Indexing API")

	service, err := gis.getIndexingService(ctx, serviceAccount)
	if err != nil {
		return &models.IndexResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get indexing service: %v", err),
			URL:     url,
		}, err
	}

	urlNotification := &indexing.UrlNotification{
		Url:  url,
		Type: "URL_UPDATED",
	}

	call := service.UrlNotifications.Publish(urlNotification)
	resp, err := call.Do()
	if err != nil {
		gis.logger.WithError(err).WithField("url", url).Error("Failed to submit URL")
		return &models.IndexResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to submit URL: %v", err),
			URL:     url,
		}, err
	}

	gis.logger.WithField("url", url).WithField("response", resp).Info("URL submitted successfully")

	return &models.IndexResponse{
		Success: true,
		Message: "URL submitted successfully",
		URL:     url,
	}, nil
}

func (gis *GoogleIndexingService) SubmitURLsBatch(ctx context.Context, urls []string, serviceAccount *models.ServiceAccountCredentials) (*models.BatchIndexResponse, error) {
	gis.logger.WithField("count", len(urls)).Info("Submitting batch URLs to Google Indexing API")

	var wg sync.WaitGroup
	results := make([]models.IndexResponse, len(urls))

	// Use goroutines for concurrent processing
	for i, url := range urls {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()

			result, err := gis.SubmitURL(ctx, u, serviceAccount)
			if err != nil {
				results[index] = models.IndexResponse{
					Success: false,
					Message: fmt.Sprintf("Failed to submit URL: %v", err),
					URL:     u,
				}
			} else {
				results[index] = *result
			}
		}(i, url)
	}

	wg.Wait()

	// Calculate statistics
	stats := models.BatchIndexResponseStats{
		Total: len(urls),
	}

	for _, result := range results {
		if result.Success {
			stats.Successful++
		} else {
			stats.Failed++
		}
	}

	response := &models.BatchIndexResponse{
		Success:    stats.Failed == 0,
		Message:    fmt.Sprintf("Processed %d URLs: %d successful, %d failed", stats.Total, stats.Successful, stats.Failed),
		Results:    results,
		Statistics: stats,
	}

	gis.logger.WithField("statistics", stats).Info("Batch URL submission completed")

	return response, nil
}

func (gis *GoogleIndexingService) GetURLStatus(ctx context.Context, url string, serviceAccount *models.ServiceAccountCredentials) (*models.StatusResponse, error) {
	gis.logger.WithField("url", url).Info("Getting URL status from Google Indexing API")

	service, err := gis.getIndexingService(ctx, serviceAccount)
	if err != nil {
		return &models.StatusResponse{
			URL:    url,
			Status: "error",
		}, err
	}

	call := service.UrlNotifications.GetMetadata()
	call.Url(url)

	resp, err := call.Do()
	if err != nil {
		gis.logger.WithError(err).WithField("url", url).Error("Failed to get URL status")
		return &models.StatusResponse{
			URL:    url,
			Status: "error",
		}, err
	}

	status := "unknown"
	lastUpdated := ""

	if resp.LatestUpdate != nil {
		status = resp.LatestUpdate.Type
		lastUpdated = resp.LatestUpdate.NotifyTime
	}

	return &models.StatusResponse{
		URL:         url,
		Status:      status,
		LastUpdated: lastUpdated,
	}, nil
}

// ClearCache clears the service cache (useful for cleanup)
func (gis *GoogleIndexingService) ClearCache() {
	gis.cacheMutex.Lock()
	defer gis.cacheMutex.Unlock()

	gis.serviceCache = make(map[string]*indexing.Service)
	gis.logger.Info("Service cache cleared")
}

// GetCacheStats returns cache statistics
func (gis *GoogleIndexingService) GetCacheStats() map[string]interface{} {
	gis.cacheMutex.RLock()
	defer gis.cacheMutex.RUnlock()

	return map[string]interface{}{
		"cached_services": len(gis.serviceCache),
		"has_default":     false, // No default service account
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
	}
}
