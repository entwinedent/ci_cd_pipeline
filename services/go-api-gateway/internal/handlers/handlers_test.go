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

func TestRobotsTxt(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()

	handler.RobotsTxt(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("Expected Content-Type text/plain, got %s", contentType)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty body")
	}
	if body != "User-agent: *\nDisallow: /api/\n" {
		t.Errorf("Unexpected robots.txt content: %s", body)
	}
}

func TestSitemap(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	w := httptest.NewRecorder()

	handler.Sitemap(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/xml" {
		t.Errorf("Expected Content-Type application/xml, got %s", contentType)
	}

	body := w.Body.String()
	if body == "" {
		t.Error("Expected non-empty body")
	}
	if !contains(body, "<?xml version") {
		t.Error("Expected XML content")
	}
	if !contains(body, "urlset") {
		t.Error("Expected urlset in sitemap")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
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
