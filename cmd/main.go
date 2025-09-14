package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"google-indexing-api/internal/config"
	"google-indexing-api/internal/handlers"
	"google-indexing-api/internal/middleware"
	"google-indexing-api/internal/services"
	"google-indexing-api/pkg/utils"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		logrus.Fatal("Failed to load configuration: ", err)
	}

	cfg := config.GetConfig()

	// Setup logger
	logger := utils.SetupLogger(cfg.Logging.Level, cfg.Logging.Format)

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Initialize Google Indexing Service
	indexingService, err := services.NewGoogleIndexingService(logger)
	if err != nil {
		logger.Fatal("Failed to initialize Google Indexing Service: ", err)
	}

	// Initialize handlers
	indexingHandler := handlers.NewIndexingHandler(indexingService, logger)

	// Setup router
	router := setupRouter(indexingHandler, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.WithField("port", cfg.Server.Port).Info("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exited")
}

func setupRouter(indexingHandler *handlers.IndexingHandler, logger *logrus.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.ErrorHandler(logger))
	router.Use(gin.Recovery())

	// Health check endpoint (publicly accessible)
	router.GET("/api/health", indexingHandler.HealthCheck)

	// API routes (no authentication required)
	api := router.Group("/api/v1")
	{
		// Single URL indexing
		api.POST("/index", indexingHandler.SubmitURL)

		// Batch URL indexing
		api.POST("/index/batch", indexingHandler.SubmitURLsBatch)

		// URL status check
		api.GET("/status/*url", func(c *gin.Context) {
			// Extract URL from path parameter
			url := c.Param("url")[1:] // Remove leading slash
			c.Params = gin.Params{{Key: "url", Value: url}}
			indexingHandler.GetURLStatus(c)
		})

		// Cache management
		api.GET("/cache/stats", indexingHandler.GetCacheStats)
		api.POST("/cache/clear", indexingHandler.ClearCache)
	}

	return router
}
