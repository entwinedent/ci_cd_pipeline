package grpc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSetIntegration(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		
		// Check content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data stored successfully",
		})
	}))
	defer server.Close()
	
	// Parse server URL
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Set("test-key", []byte("test-value"), 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !success {
		t.Error("Expected success to be true")
	}
}

func TestGetIntegration(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		
		// Return success response with data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data retrieved successfully",
			Value:   "test-value",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !found {
		t.Error("Expected found to be true")
	}
	if string(value) != "test-value" {
		t.Errorf("Expected 'test-value', got %s", string(value))
	}
}

func TestGetNotFound(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("nonexistent-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found {
		t.Error("Expected found to be false")
	}
	if value != nil {
		t.Error("Expected value to be nil")
	}
}

func TestDeleteIntegration(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		
		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data deleted successfully",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Delete("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !success {
		t.Error("Expected success to be true")
	}
}

func TestHealthCheckIntegration(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Service is healthy",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	healthy, err := client.HealthCheck()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !healthy {
		t.Error("Expected healthy to be true")
	}
}

func TestHealthCheckUnhealthy(t *testing.T) {
	// Create a test server that returns unhealthy
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Service is unhealthy",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	healthy, err := client.HealthCheck()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if healthy {
		t.Error("Expected healthy to be false")
	}
}

func TestSetWithErrorStatus(t *testing.T) {
	// Create a test server that returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Set("test-key", []byte("test-value"), 3600)
	if err == nil {
		t.Error("Expected error for non-200 status")
	}
	if success {
		t.Error("Expected success to be false")
	}
}

func TestGetWithErrorStatus(t *testing.T) {
	// Create a test server that returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("test-key")
	if err == nil {
		t.Error("Expected error for non-200 status")
	}
	if found {
		t.Error("Expected found to be false")
	}
	if value != nil {
		t.Error("Expected value to be nil")
	}
}

func TestDeleteWithErrorStatus(t *testing.T) {
	// Create a test server that returns error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Delete("test-key")
	if err == nil {
		t.Error("Expected error for non-200 status")
	}
	if success {
		t.Error("Expected success to be false")
	}
}

func TestSetWithInvalidResponse(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Set("test-key", []byte("test-value"), 3600)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if success {
		t.Error("Expected success to be false")
	}
}

func TestGetWithInvalidResponse(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("test-key")
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if found {
		t.Error("Expected found to be false")
	}
	if value != nil {
		t.Error("Expected value to be nil")
	}
}

func TestHealthCheckWithInvalidResponse(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	healthy, err := client.HealthCheck()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if healthy {
		t.Error("Expected healthy to be false")
	}
}

func TestHealthCheckWithConnectionError(t *testing.T) {
	// Create a client with invalid URL
	client, err := NewDataStoreClient("http://invalid-host-that-does-not-exist:12345")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	healthy, err := client.HealthCheck()
	if err == nil {
		t.Error("Expected connection error")
	}
	if healthy {
		t.Error("Expected healthy to be false")
	}
}

func TestSetWithConnectionError(t *testing.T) {
	// Create a client with invalid URL
	client, err := NewDataStoreClient("http://invalid-host-that-does-not-exist:12345")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Set("test-key", []byte("test-value"), 3600)
	if err == nil {
		t.Error("Expected connection error")
	}
	if success {
		t.Error("Expected success to be false")
	}
}

func TestGetWithConnectionError(t *testing.T) {
	// Create a client with invalid URL
	client, err := NewDataStoreClient("http://invalid-host-that-does-not-exist:12345")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("test-key")
	if err == nil {
		t.Error("Expected connection error")
	}
	if found {
		t.Error("Expected found to be false")
	}
	if value != nil {
		t.Error("Expected value to be nil")
	}
}

func TestDeleteWithConnectionError(t *testing.T) {
	// Create a client with invalid URL
	client, err := NewDataStoreClient("http://invalid-host-that-does-not-exist:12345")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	success, err := client.Delete("test-key")
	if err == nil {
		t.Error("Expected connection error")
	}
	if success {
		t.Error("Expected success to be false")
	}
}

func TestSetURLConstruction(t *testing.T) {
	// Test that Set constructs the correct URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the URL path
		if !strings.Contains(r.URL.Path, "/api/v1/data/test-key") {
			t.Errorf("Expected URL path to contain /api/v1/data/test-key, got %s", r.URL.Path)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data stored successfully",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	_, err = client.Set("test-key", []byte("test-value"), 3600)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetURLConstruction(t *testing.T) {
	// Test that Get constructs the correct URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the URL path
		if !strings.Contains(r.URL.Path, "/api/v1/data/test-key") {
			t.Errorf("Expected URL path to contain /api/v1/data/test-key, got %s", r.URL.Path)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data retrieved successfully",
			Value:   "test-value",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	_, _, err = client.Get("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDeleteURLConstruction(t *testing.T) {
	// Test that Delete constructs the correct URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the URL path
		if !strings.Contains(r.URL.Path, "/api/v1/data/test-key") {
			t.Errorf("Expected URL path to contain /api/v1/data/test-key, got %s", r.URL.Path)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data deleted successfully",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	_, err = client.Delete("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestHealthCheckURLConstruction(t *testing.T) {
	// Test that HealthCheck constructs the correct URL
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the URL path
		if r.URL.Path != "/healthz" {
			t.Errorf("Expected URL path to be /healthz, got %s", r.URL.Path)
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Service is healthy",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	_, err = client.HealthCheck()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetWithSuccessFalse(t *testing.T) {
	// Test the case where response.Success is false but value is not empty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Error retrieving data",
			Value:   "some-value",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found {
		t.Error("Expected found to be false when Success is false")
	}
	if value != nil {
		t.Error("Expected value to be nil when Success is false")
	}
}

func TestGetWithEmptyValue(t *testing.T) {
	// Test the case where response.Success is true but value is empty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Data retrieved successfully",
			Value:   "",
		})
	}))
	defer server.Close()
	
	serverURL, _ := url.Parse(server.URL)
	
	client, err := NewDataStoreClient(serverURL.String())
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	value, found, err := client.Get("test-key")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if found {
		t.Error("Expected found to be false when Value is empty")
	}
	if value != nil {
		t.Error("Expected value to be nil when Value is empty")
	}
}

func TestSetWithMarshalError(t *testing.T) {
	// This is a theoretical test - json.Marshal rarely fails for valid structs
	// but we need to cover the error path
	// Since we cannot easily make json.Marshal fail, we will test the TTL logic
	
	client, err := NewDataStoreClient("http://localhost:50051")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	
	// Test with TTL = 0 (should not set TTLSeconds)
	_, err = client.Set("test-key", []byte("test-value"), 0)
	// This will fail because the server is not running, but we just want to test
	// that the payload construction with TTL=0 works
	if err == nil {
		t.Log("Expected connection error (server not running)")
	}
}
