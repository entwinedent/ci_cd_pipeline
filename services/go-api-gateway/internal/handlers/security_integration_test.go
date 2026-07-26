package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityEndpointsIntegration(t *testing.T) {
	// Create a mock data store client
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	// Test robots.txt endpoint
	t.Run("RobotsTxt", func(t *testing.T) {
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
		expected := "User-agent: *\nDisallow: /api/\n"
		if body != expected {
			t.Errorf("Expected robots.txt content: %s, got: %s", expected, body)
		}
	})

	// Test sitemap.xml endpoint
	t.Run("Sitemap", func(t *testing.T) {
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
		if len(body) == 0 {
			t.Error("Expected non-empty sitemap body")
		}
	})
}
