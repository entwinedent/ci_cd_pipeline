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
	dataStoreTarget := getEnv("DATA_STORE_TARGET", "127.0.0.1:50051")
	// Ensure data store target has http:// prefix if not already present
	if dataStoreTarget != "" && !hasProtocol(dataStoreTarget) {
		dataStoreTarget = "http://" + dataStoreTarget
	}
	return &Config{
		Port:            getEnv("PORT", "8080"),
		DataStoreTarget: dataStoreTarget,
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		FeatureFlagURL:  getEnv("FEATURE_FLAG_URL", ""),
	}
}

func hasProtocol(url string) bool {
	return len(url) >= 7 && (url[:7] == "http://" || url[:8] == "https://")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
