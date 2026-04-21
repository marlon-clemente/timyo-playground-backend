package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the configuration for the application.
type Config struct {
	Port        string
	ServiceName string
	Development bool

	JWTSecret       string
	FrontendBaseUrl string
	AuthServiceURL string

	LogFormat string

	OtelEndpoint string
	OtelAuth     string
	OtelInsecure bool
	OtelOrg      string

	DatabaseDSN string
}

// Load loads the configuration from environment variables.
// It attempts to load from a .env file first, but doesn't fail if the file is missing.
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	return &Config{
		Development:             os.Getenv("DEVELOPMENT") == "true",
		JWTSecret:               getEnv("JWT_SECRET", ""),
		Port:                    getEnv("PORT", "3003"),
		ServiceName:             getEnv("SERVICE_NAME", "playground-backend"),
		LogFormat:               getEnv("LOG_FORMAT", "text"),
		OtelEndpoint:            getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", ""),
		OtelAuth:                getEnv("OTEL_EXPORTER_OTLP_AUTH", ""),
		OtelInsecure:            os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true",
		OtelOrg:                 getEnv("OTEL_OPENOBSERVE_ORG", "default"),
		DatabaseDSN:             getEnv("DATABASE_DSN", ""),
		FrontendBaseUrl:         getEnv("FRONTEND_BASE_URL", ""),
		AuthServiceURL:          getEnv("AUTH_SERVICE_URL", ""),
	}
}

// getEnv retrieves an environment variable or returns a default value if not set.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
