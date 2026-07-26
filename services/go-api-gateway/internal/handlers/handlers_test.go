package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSmoke(t *testing.T) {
	// Smoke test to verify basic functionality
	handler := NewHTTPHandler(nil)
	if handler == nil {
		t.Error("Handler creation failed")
	}
}

func TestHealthCheck(t *testing.T) {
	// Create a mock data store client that always returns healthy
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestLivenessCheck(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/livez", nil)
	w := httptest.NewRecorder()

	handler.LivenessCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// MockDataStoreClient is a mock implementation for testing
type MockDataStoreClient struct {
	healthy bool
}

func (m *MockDataStoreClient) Set(key string, value []byte, ttl int64) (bool, error) {
	return true, nil
}

func (m *MockDataStoreClient) Get(key string) ([]byte, bool, error) {
	return []byte("test"), true, nil
}

func (m *MockDataStoreClient) Delete(key string) (bool, error) {
	return true, nil
}

func (m *MockDataStoreClient) HealthCheck() (bool, error) {
	return m.healthy, nil
}
