package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	Port            string
	DataStoreTarget string
	LogLevel        string
	FeatureFlagURL  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DataStoreTarget: getEnv("DATA_STORE_TARGET", "localhost:50051"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		FeatureFlagURL:  getEnv("FEATURE_FLAG_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
