//go:build integration
// +build integration

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/config"
	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/grpc"
)

func TestMainIntegration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Get data store target from environment
	dataStoreTarget := os.Getenv("DATA_STORE_TARGET")
	if dataStoreTarget == "" {
		dataStoreTarget = "http://localhost:50051"
	}

	// Wait for mock server to be ready
	client := &http.Client{Timeout: 5 * time.Second}
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(dataStoreTarget + "/healthz")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			break
		}
		if i == maxRetries-1 {
			t.Fatalf("Mock server not ready after %d attempts", maxRetries)
		}
		time.Sleep(1 * time.Second)
	}

	// Test gRPC client connection
	dataStoreClient, err := grpc.NewDataStoreClient(dataStoreTarget)
	if err != nil {
		t.Fatalf("Failed to create data store client: %v", err)
	}
	defer dataStoreClient.Close()

	// Test health check
	healthy, err := dataStoreClient.HealthCheck()
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
	if !healthy {
		t.Error("Expected healthy to be true")
	}

	// Test Set operation
	success, err := dataStoreClient.Set("test-key", []byte("test-value"), 3600)
	if err != nil {
		t.Errorf("Set operation failed: %v", err)
	}
	if !success {
		t.Error("Expected Set to succeed")
	}

	// Test Get operation
	value, found, err := dataStoreClient.Get("test-key")
	if err != nil {
		t.Errorf("Get operation failed: %v", err)
	}
	if !found {
		t.Error("Expected key to be found")
	}
	if string(value) != "test-value" {
		t.Errorf("Expected 'test-value', got %s", string(value))
	}

	// Test Delete operation
	success, err = dataStoreClient.Delete("test-key")
	if err != nil {
		t.Errorf("Delete operation failed: %v", err)
	}
	if !success {
		t.Error("Expected Delete to succeed")
	}

	// Test server creation
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: dataStoreTarget,
		LogLevel:        "info",
	}

	server := NewServer(cfg, dataStoreClient)
	if server == nil {
		t.Fatal("Expected server to be created")
	}

	// Test server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	_ = server.httpServer.Shutdown(ctx)
}

func TestServerWithRealDataStore(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	dataStoreTarget := os.Getenv("DATA_STORE_TARGET")
	if dataStoreTarget == "" {
		dataStoreTarget = "http://localhost:50051"
	}

	// Wait for server
	client := &http.Client{Timeout: 5 * time.Second}
	for i := 0; i < 10; i++ {
		resp, err := client.Get(dataStoreTarget + "/healthz")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			break
		}
		if i == 9 {
			t.Fatalf("Mock server not ready")
		}
		time.Sleep(500 * time.Millisecond)
	}

	dataStoreClient, err := grpc.NewDataStoreClient(dataStoreTarget)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer dataStoreClient.Close()

	cfg := &config.Config{
		Port:            "8081",
		DataStoreTarget: dataStoreTarget,
		LogLevel:        "debug",
	}

	server := NewServer(cfg, dataStoreClient)
	if server == nil {
		t.Fatal("Expected server to be created")
	}

	// Start server
	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(200 * time.Millisecond)

	// Test health endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/healthz", cfg.Port))
	if err != nil {
		t.Errorf("Failed to call health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.httpServer.Shutdown(ctx)
}
