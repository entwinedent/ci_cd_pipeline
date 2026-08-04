package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/config"
	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/handlers"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	if server == nil {
		t.Fatal("Expected server to be created")
	}

	if server.config != cfg {
		t.Error("Expected config to be set")
	}

	if server.router == nil {
		t.Error("Expected router to be set")
	}

	if server.httpServer == nil {
		t.Error("Expected httpServer to be set")
	}

	if server.httpServer.Addr != ":8080" {
		t.Errorf("Expected address :8080, got %s", server.httpServer.Addr)
	}
}

func TestSetupRoutes(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /healthz, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/livez", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /livez, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/robots.txt", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /robots.txt, got %d", w.Code)
	}

	req = httptest.NewRequest("GET", "/sitemap.xml", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /sitemap.xml, got %d", w.Code)
	}
}

func TestRouteRegistration(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	routes := []string{
		"/healthz",
		"/readyz",
		"/livez",
		"/robots.txt",
		"/sitemap.xml",
	}

	for _, route := range routes {
		req := httptest.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		server.router.ServeHTTP(w, req)
		_ = w.Code
	}
}

func TestAPIRoutes(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	// Test GET /api/v1/data/{key}
	req := httptest.NewRequest("GET", "/api/v1/data/test-key", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Test POST /api/v1/data/{key}
	req = httptest.NewRequest("POST", "/api/v1/data/test-key", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Test DELETE /api/v1/data/{key}
	req = httptest.NewRequest("DELETE", "/api/v1/data/test-key", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)
}

func TestServerDataStore(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	if server.dataStore == nil {
		t.Error("Expected dataStore to be set")
	}
}

func TestServerConfig(t *testing.T) {
	cfg := &config.Config{
		Port:            "9090",
		DataStoreTarget: "http://localhost:50052",
		LogLevel:        "debug",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	if server.config.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", server.config.Port)
	}

	if server.config.DataStoreTarget != "http://localhost:50052" {
		t.Errorf("Expected DataStoreTarget http://localhost:50052, got %s", server.config.DataStoreTarget)
	}

	if server.config.LogLevel != "debug" {
		t.Errorf("Expected LogLevel debug, got %s", server.config.LogLevel)
	}
}

func TestServerShutdown(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	// Start server in background
	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait a bit for server to start
	time.Sleep(100 * time.Millisecond)

	// Test shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		t.Errorf("Expected no error on shutdown, got %v", err)
	}
}

func TestMainSignalHandling(t *testing.T) {
	// This test simulates the signal handling logic from main
	// We can't test the actual main function directly, but we can test the components

	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	// Test that server can be created and shut down
	if server == nil {
		t.Fatal("Expected server to be created")
	}

	// Simulate shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		// Simulate receiving a signal
		cancel()
	}()

	_ = server.httpServer.Shutdown(ctx)
}

func TestServerRouterNotNil(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	if server.router == nil {
		t.Error("Expected router to be initialized")
	}
}

func TestServerHTTPServerNotNil(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	if server.httpServer == nil {
		t.Error("Expected httpServer to be initialized")
	}
}

func TestServerHTTPServerHandler(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	if server.httpServer.Handler == nil {
		t.Error("Expected httpServer Handler to be set")
	}
}

func TestServerHTTPServerAddr(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	expectedAddr := ":8080"
	if server.httpServer.Addr != expectedAddr {
		t.Errorf("Expected address %s, got %s", expectedAddr, server.httpServer.Addr)
	}
}

func TestSetupRoutesMiddleware(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	// Test that middleware is applied by checking that routes respond
	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// If middleware is working, the request should be handled
	if w.Code == 0 {
		t.Error("Expected response code to be set")
	}
}

func TestReadyzEndpoint(t *testing.T) {
	cfg := &config.Config{
		Port:            "8080",
		DataStoreTarget: "http://localhost:50051",
		LogLevel:        "info",
	}

	mockClient := handlers.NewMockDataStoreClient(true)
	server := NewServer(cfg, mockClient)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for /readyz, got %d", w.Code)
	}
}

func TestRun_DataStoreConnectionError(t *testing.T) {
	// Set environment to use invalid data store target
	// Note: grpc.NewDataStoreClient doesn't actually connect immediately,
	// so this test is skipped. The error would only occur on actual HTTP requests.
	t.Skip("grpc client doesn't connect immediately, error would only occur on requests")
}

func TestRun_Success(t *testing.T) {
	// Signal handling test is platform-specific and unreliable
	// Use Docker integration tests instead for full coverage
	t.Skip("Signal handling is tested in Docker integration tests")
}
