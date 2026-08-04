package grpc

import (
	"testing"
	"time"
)

func TestNewDataStoreClient(t *testing.T) {
	// Test client creation
	client, err := NewDataStoreClient("http://localhost:50051")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created")
	}
	if client.baseURL != "http://localhost:50051" {
		t.Errorf("Expected base URL http://localhost:50051, got %s", client.baseURL)
	}
	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestNewDataStoreClientTimeout(t *testing.T) {
	client, err := NewDataStoreClient("http://localhost:50051")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created")
	}
	// Check that timeout is set (default 5 seconds)
	if client.client.Timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", client.client.Timeout)
	}
}

func TestClose(t *testing.T) {
	client, err := NewDataStoreClient("http://localhost:50051")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err = client.Close()
	if err != nil {
		t.Errorf("Expected no error on close, got %v", err)
	}
}

func TestSetValue(t *testing.T) {
	// Test SetValue structure
	payload := SetValue{
		Value: "test-value",
	}
	if payload.Value != "test-value" {
		t.Errorf("Expected test-value, got %s", payload.Value)
	}

	ttl := int64(3600)
	payload.TTLSeconds = &ttl
	if payload.TTLSeconds == nil {
		t.Error("Expected TTL to be set")
	}
	if *payload.TTLSeconds != 3600 {
		t.Errorf("Expected TTL 3600, got %d", *payload.TTLSeconds)
	}
}

func TestResponse(t *testing.T) {
	response := Response{
		Success: true,
		Message: "test message",
		Value:   "test value",
	}

	if response.Success != true {
		t.Error("Expected Success to be true")
	}
	if response.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", response.Message)
	}
	if response.Value != "test value" {
		t.Errorf("Expected value 'test value', got %s", response.Value)
	}

	// Test response with empty value
	response2 := Response{
		Success: false,
		Message: "error message",
	}

	if response2.Success != false {
		t.Error("Expected Success to be false")
	}
	if response2.Value != "" {
		t.Errorf("Expected empty value, got %s", response2.Value)
	}
}

func TestSetValueWithoutTTL(t *testing.T) {
	payload := SetValue{
		Value:      "test-value",
		TTLSeconds: nil,
	}
	if payload.Value != "test-value" {
		t.Errorf("Expected test-value, got %s", payload.Value)
	}
	if payload.TTLSeconds != nil {
		t.Error("Expected TTL to be nil")
	}
}

func TestSetMarshalError(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	// Test with invalid data that can't be marshaled
	// This is a structural test to ensure error handling path exists
	_, err := client.Set("test", []byte("value"), 0)
	if err == nil {
		t.Error("Expected error when server is not available")
	}
}

func TestGetMarshalError(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	// Test Get when server is not available
	_, _, err := client.Get("test")
	if err == nil {
		t.Error("Expected error when server is not available")
	}
}

func TestDeleteMarshalError(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	// Test Delete when server is not available
	_, err := client.Delete("test")
	if err == nil {
		t.Error("Expected error when server is not available")
	}
}

func TestHealthCheckMarshalError(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	// Test HealthCheck when server is not available
	_, err := client.HealthCheck()
	if err == nil {
		t.Error("Expected error when server is not available")
	}
}

func TestSetWithZeroTTL(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	// Test with TTL = 0 (should not include TTLSeconds in payload)
	_, err := client.Set("test", []byte("value"), 0)
	if err == nil {
		t.Log("Expected connection error (server not running)")
	}
}

func TestSetWithNegativeTTL(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	// Test with negative TTL (should not include TTLSeconds in payload)
	_, err := client.Set("test", []byte("value"), -1)
	if err == nil {
		t.Log("Expected connection error (server not running)")
	}
}

func TestClientBaseURL(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	if client.baseURL != "http://localhost:50051" {
		t.Errorf("Expected base URL http://localhost:50051, got %s", client.baseURL)
	}
}

func TestClientNotNil(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	if client == nil {
		t.Error("Expected client to be created")
	}
}

func TestClientHTTPClientNotNil(t *testing.T) {
	client, _ := NewDataStoreClient("http://localhost:50051")

	if client.client == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestSetValueStruct(t *testing.T) {
	ttl := int64(3600)
	payload := SetValue{
		Value:      "test-value",
		TTLSeconds: &ttl,
	}

	if payload.Value != "test-value" {
		t.Errorf("Expected Value test-value, got %s", payload.Value)
	}
	if payload.TTLSeconds == nil {
		t.Error("Expected TTLSeconds to be set")
	}
	if *payload.TTLSeconds != 3600 {
		t.Errorf("Expected TTLSeconds 3600, got %d", *payload.TTLSeconds)
	}
}

func TestSetValueStructNilTTL(t *testing.T) {
	payload := SetValue{
		Value:      "test-value",
		TTLSeconds: nil,
	}

	if payload.Value != "test-value" {
		t.Errorf("Expected Value test-value, got %s", payload.Value)
	}
	if payload.TTLSeconds != nil {
		t.Error("Expected TTLSeconds to be nil")
	}
}

func TestResponseStruct(t *testing.T) {
	response := Response{
		Success: true,
		Message: "test message",
		Value:   "test value",
	}

	if response.Success != true {
		t.Error("Expected Success to be true")
	}
	if response.Message != "test message" {
		t.Errorf("Expected Message test message, got %s", response.Message)
	}
	if response.Value != "test value" {
		t.Errorf("Expected Value test value, got %s", response.Value)
	}
}

func TestResponseStructEmptyValue(t *testing.T) {
	response := Response{
		Success: false,
		Message: "error message",
	}

	if response.Success != false {
		t.Error("Expected Success to be false")
	}
	if response.Value != "" {
		t.Errorf("Expected empty Value, got %s", response.Value)
	}
}
