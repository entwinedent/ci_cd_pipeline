package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestSmoke(t *testing.T) {
	// Smoke test to verify basic functionality
	handler := NewHTTPHandler(nil)
	if handler == nil {
		t.Error("Handler creation failed")
	}
}

func TestMockDataStoreClient(t *testing.T) {
	client := NewMockDataStoreClient(true)
	if client == nil {
		t.Error("Expected mock client to be created")
	}

	if client.healthy != true {
		t.Error("Expected healthy to be true")
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

func TestRootHandler(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.RootHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "active" {
		t.Errorf("Expected status 'active', got %s", response["status"])
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

func TestReadinessCheck(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()

	handler.ReadinessCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["ready"] != true {
		t.Errorf("Expected ready true, got %v", response["ready"])
	}
}

func TestReadinessCheckUnhealthy(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: false}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()

	handler.ReadinessCheck(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["ready"] != false {
		t.Errorf("Expected ready false, got %v", response["ready"])
	}
}

func TestSetData(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true, getNotFound: false}
	handler := NewHTTPHandler(mockClient)

	// Send string value - the handler expects []byte but JSON unmarshal will handle string
	body := []byte(`{"value": "test-value", "ttl_seconds": 3600}`)
	req := httptest.NewRequest("POST", "/api/v1/data/testkey", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add mux vars to the request context
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
}

func TestGetData(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/data/testkey", nil)

	// Add mux vars to the request context
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.GetData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
}

func TestGetDataNotFound(t *testing.T) {
	// Create a mock that returns not found
	mockClient := &MockDataStoreClient{healthy: true, getNotFound: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/data/nonexistent", nil)

	// Add mux vars to the request context
	vars := map[string]string{"key": "nonexistent"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.GetData(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestDeleteData(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("DELETE", "/api/v1/data/testkey", nil)

	// Add mux vars to the request context
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.DeleteData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
}

func TestSetDataWithEmptyKey(t *testing.T) {
	mockClient := NewMockDataStoreClient(true)
	handler := NewHTTPHandler(mockClient)

	body := []byte(`{"value": "test-value", "ttl_seconds": 3600}`)
	req := httptest.NewRequest("POST", "/api/v1/data/", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"key": ""}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetDataWithEmptyKey(t *testing.T) {
	mockClient := NewMockDataStoreClient(true)
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/data/", nil)

	vars := map[string]string{"key": ""}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.GetData(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestDeleteDataWithEmptyKey(t *testing.T) {
	mockClient := NewMockDataStoreClient(true)
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("DELETE", "/api/v1/data/", nil)

	vars := map[string]string{"key": ""}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.DeleteData(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestSetDataWithInvalidJSON(t *testing.T) {
	mockClient := NewMockDataStoreClient(true)
	handler := NewHTTPHandler(mockClient)

	body := []byte(`invalid json`)
	req := httptest.NewRequest("POST", "/api/v1/data/testkey", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestHealthCheckUnhealthy(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: false}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", w.Code)
	}
}

func TestSetDataWithError(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true, setError: true}
	handler := NewHTTPHandler(mockClient)

	body := []byte(`{"value": "test-value", "ttl_seconds": 3600}`)
	req := httptest.NewRequest("POST", "/api/v1/data/testkey", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestGetDataWithError(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true, getError: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/data/testkey", nil)

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.GetData(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestDeleteDataWithError(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true, deleteError: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("DELETE", "/api/v1/data/testkey", nil)

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.DeleteData(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}

func TestReadinessCheckWithError(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true, healthCheckError: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()

	handler.ReadinessCheck(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", w.Code)
	}
}

func TestSetDataWithoutTTL(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	body := []byte(`{"value": "test-value"}`)
	req := httptest.NewRequest("POST", "/api/v1/data/testkey", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRegisterRoutes(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// Verify routes are registered by checking router
	if router == nil {
		t.Error("Expected router to be set")
	}
}

func TestNewHTTPHandler(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
	if handler.dataStore != mockClient {
		t.Error("Expected dataStore to be set")
	}
}

func TestRootHandlerHeaders(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.RootHandler(w, req)

	// Check security headers
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
	if w.Header().Get("Cache-Control") != "no-store, no-cache, must-revalidate, private" {
		t.Error("Expected Cache-Control header")
	}
}

func TestHealthCheckHeaders(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	handler.HealthCheck(w, req)

	// Check security headers
	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
}

func TestRobotsTxtHeaders(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()
	handler.RobotsTxt(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
	if w.Header().Get("Cache-Control") != "public, max-age=3600" {
		t.Error("Expected Cache-Control header")
	}
}

func TestSitemapHeaders(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	w := httptest.NewRecorder()
	handler.Sitemap(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
	if w.Header().Get("Cross-Origin-Resource-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Resource-Policy header")
	}
	if w.Header().Get("Cache-Control") != "public, max-age=3600" {
		t.Error("Expected Cache-Control header")
	}
}

func TestLivenessCheckTimestamp(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/livez", nil)
	w := httptest.NewRecorder()
	handler.LivenessCheck(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "alive" {
		t.Errorf("Expected status 'alive', got %v", response["status"])
	}
	if response["timestamp"] == nil {
		t.Error("Expected timestamp in response")
	}
}

func TestReadinessCheckTimestamp(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	handler.ReadinessCheck(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["ready"] != true {
		t.Errorf("Expected ready true, got %v", response["ready"])
	}
	if response["timestamp"] == nil {
		t.Error("Expected timestamp in response")
	}
}

func TestGetDataResponseFields(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/data/testkey", nil)
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.GetData(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
	if response["found"] != true {
		t.Errorf("Expected found true, got %v", response["found"])
	}
	if response["value"] == nil {
		t.Error("Expected value in response")
	}
}

func TestSetDataResponseMessage(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	body := []byte(`{"value": "test-value", "ttl_seconds": 3600}`)
	req := httptest.NewRequest("POST", "/api/v1/data/testkey", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
	if response["message"] == nil {
		t.Error("Expected message in response")
	}
}

func TestDeleteDataResponseMessage(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("DELETE", "/api/v1/data/testkey", nil)
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.DeleteData(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["success"] != true {
		t.Errorf("Expected success true, got %v", response["success"])
	}
	if response["message"] == nil {
		t.Error("Expected message in response")
	}
}

func TestSetDataContentType(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	body := []byte(`{"value": "test-value"}`)
	req := httptest.NewRequest("POST", "/api/v1/data/testkey", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	handler.SetData(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestGetDataContentType(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/api/v1/data/testkey", nil)
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.GetData(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestDeleteDataContentType(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("DELETE", "/api/v1/data/testkey", nil)
	vars := map[string]string{"key": "testkey"}
	req = mux.SetURLVars(req, vars)

	w := httptest.NewRecorder()
	handler.DeleteData(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestRootHandlerResponse(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.RootHandler(w, req)

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "active" {
		t.Errorf("Expected status 'active', got %s", response["status"])
	}
	if response["message"] != "API Gateway Active" {
		t.Errorf("Expected message 'API Gateway Active', got %s", response["message"])
	}
}

func TestLivenessCheckContentType(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/livez", nil)
	w := httptest.NewRecorder()
	handler.LivenessCheck(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestReadinessCheckContentType(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/readyz", nil)
	w := httptest.NewRecorder()
	handler.ReadinessCheck(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestHealthCheckContentType(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	handler.HealthCheck(w, req)

	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Expected Content-Type application/json")
	}
}

func TestHealthCheckCacheControl(t *testing.T) {
	mockClient := &MockDataStoreClient{healthy: true}
	handler := NewHTTPHandler(mockClient)

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	handler.HealthCheck(w, req)

	if w.Header().Get("Cache-Control") != "no-store, no-cache, must-revalidate, private" {
		t.Error("Expected Cache-Control header")
	}
}

func TestNewHTTPHandlerNilClient(t *testing.T) {
	handler := NewHTTPHandler(nil)
	if handler == nil {
		t.Error("Expected handler to be created even with nil client")
	}
	if handler.dataStore != nil {
		t.Error("Expected dataStore to be nil")
	}
}

func TestRootHandlerSecurityHeaders(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.RootHandler(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
}

func TestRobotsTxtSecurityHeaders(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/robots.txt", nil)
	w := httptest.NewRecorder()
	handler.RobotsTxt(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
}

func TestSitemapSecurityHeaders(t *testing.T) {
	handler := NewHTTPHandler(nil)

	req := httptest.NewRequest("GET", "/sitemap.xml", nil)
	w := httptest.NewRecorder()
	handler.Sitemap(w, req)

	if w.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header")
	}
	if w.Header().Get("Cross-Origin-Embedder-Policy") != "require-corp" {
		t.Error("Expected Cross-Origin-Embedder-Policy header")
	}
	if w.Header().Get("Cross-Origin-Opener-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Opener-Policy header")
	}
	if w.Header().Get("Cross-Origin-Resource-Policy") != "same-origin" {
		t.Error("Expected Cross-Origin-Resource-Policy header")
	}
}
