package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server struct {
		Port    string
		GinMode string
	}
	App struct {
		Version string
		Env     string
	}
	Logging struct {
		Level            string
		Format           string
		EnableRequestLog bool
	}
	RateLimit struct {
		PerMinute int
	}
	CORS struct {
		AllowedOrigins []string
		AllowedMethods []string
		AllowedHeaders []string
	}
	Performance struct {
		CacheTTLMinutes       int
		MaxConcurrentRequests int
		RequestTimeoutSeconds int
		MaxBatchSize          int
		MaxRetryAttempts      int
		RetryDelaySeconds     int
	}
	Security struct {
		EnableSecurityHeaders bool
		TrustedProxies        []string
		EnableMetrics         bool
	}
}

var AppConfig *Config

func LoadConfig() error {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	config := &Config{}

	// Server configuration
	config.Server.Port = getEnv("PORT", "8080")
	config.Server.GinMode = getEnv("GIN_MODE", "release")

	// App configuration
	config.App.Version = getEnv("APP_VERSION", "1.0.0")
	config.App.Env = getEnv("APP_ENV", "production")

	// Logging configuration
	config.Logging.Level = getEnv("LOG_LEVEL", "info")
	config.Logging.Format = getEnv("LOG_FORMAT", "json")
	config.Logging.EnableRequestLog = getEnvBool("ENABLE_REQUEST_LOGGING", true)

	// Rate limiting configuration
	config.RateLimit.PerMinute = getEnvInt("RATE_LIMIT_PER_MINUTE", 60)

	// CORS configuration
	originsStr := getEnv("CORS_ALLOWED_ORIGINS", "*")
	if originsStr == "*" {
		config.CORS.AllowedOrigins = []string{"*"}
	} else {
		config.CORS.AllowedOrigins = strings.Split(originsStr, ",")
	}

	methodsStr := getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,OPTIONS")
	config.CORS.AllowedMethods = strings.Split(methodsStr, ",")

	headersStr := getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,X-Requested-With")
	config.CORS.AllowedHeaders = strings.Split(headersStr, ",")

	// Performance configuration
	config.Performance.CacheTTLMinutes = getEnvInt("CACHE_TTL_MINUTES", 60)
	config.Performance.MaxConcurrentRequests = getEnvInt("MAX_CONCURRENT_REQUESTS", 10)
	config.Performance.RequestTimeoutSeconds = getEnvInt("REQUEST_TIMEOUT_SECONDS", 30)
	config.Performance.MaxBatchSize = getEnvInt("MAX_BATCH_SIZE", 100)
	config.Performance.MaxRetryAttempts = getEnvInt("MAX_RETRY_ATTEMPTS", 3)
	config.Performance.RetryDelaySeconds = getEnvInt("RETRY_DELAY_SECONDS", 2)

	// Security configuration
	config.Security.EnableSecurityHeaders = getEnvBool("ENABLE_SECURITY_HEADERS", true)
	config.Security.EnableMetrics = getEnvBool("ENABLE_METRICS", true)

	trustedProxiesStr := getEnv("TRUSTED_PROXIES", "127.0.0.1")
	config.Security.TrustedProxies = strings.Split(trustedProxiesStr, ",")

	AppConfig = config
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func GetConfig() *Config {
	return AppConfig
}
