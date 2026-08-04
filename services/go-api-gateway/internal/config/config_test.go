package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Test default configuration
	config := LoadConfig()
	if config == nil {
		t.Error("Expected config to be loaded")
	}
	if config.DataStoreTarget == "" {
		t.Error("Expected DataStoreTarget to be set")
	}
	if config.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Port)
	}
}

func TestLoadConfigWithEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DATA_STORE_TARGET", "localhost:50052")
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("DATA_STORE_TARGET")
	defer os.Unsetenv("LOG_LEVEL")

	config := LoadConfig()
	if config.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", config.Port)
	}
	if config.DataStoreTarget != "http://localhost:50052" {
		t.Errorf("Expected http://localhost:50052, got %s", config.DataStoreTarget)
	}
	if config.LogLevel != "debug" {
		t.Errorf("Expected log level debug, got %s", config.LogLevel)
	}
}

func TestHasProtocol(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://localhost", true},
		{"https://example.com", true},
		{"localhost:8080", false},
		{"", false},
		{"ftp://server.com", false},
	}

	for _, test := range tests {
		result := hasProtocol(test.url)
		if result != test.expected {
			t.Errorf("hasProtocol(%s) = %v, expected %v", test.url, result, test.expected)
		}
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "default")
	if result != "test_value" {
		t.Errorf("Expected test_value, got %s", result)
	}

	// Test with non-existing env var
	result = getEnv("NON_EXISTING_VAR", "default")
	if result != "default" {
		t.Errorf("Expected default, got %s", result)
	}
}
