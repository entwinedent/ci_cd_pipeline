package grpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DataStoreClient is an HTTP client for the Data Store service
type DataStoreClient struct {
	client  *http.Client
	baseURL string
}

// NewDataStoreClient creates a new HTTP client for the Data Store service
func NewDataStoreClient(address string) (*DataStoreClient, error) {
	return &DataStoreClient{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		baseURL: address,
	}, nil
}

// Close closes the HTTP client (no-op for HTTP client)
func (c *DataStoreClient) Close() error {
	return nil
}

// SetValue represents the JSON payload for setting a value
type SetValue struct {
	Value      string `json:"value"`
	TTLSeconds *int64 `json:"ttl_seconds,omitempty"`
}

// Response represents the JSON response from the data store
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Set stores a key-value pair with optional TTL
func (c *DataStoreClient) Set(key string, value []byte, ttl int64) (bool, error) {
	payload := SetValue{
		Value: string(value),
	}
	if ttl > 0 {
		payload.TTLSeconds = &ttl
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	url := fmt.Sprintf("%s/api/v1/data/%s", c.baseURL, key)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	return response.Success, nil
}

// Get retrieves a value by key
func (c *DataStoreClient) Get(key string) ([]byte, bool, error) {
	url := fmt.Sprintf("%s/api/v1/data/%s", c.baseURL, key)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, false, err
	}

	if !response.Success || response.Value == "" {
		return nil, false, nil
	}

	return []byte(response.Value), true, nil
}

// Delete removes a key-value pair
func (c *DataStoreClient) Delete(key string) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/data/%s", c.baseURL, key)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return false, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, err
	}

	return response.Success, nil
}

// HealthCheck checks the health of the data store service
func (c *DataStoreClient) HealthCheck() (bool, error) {
	url := fmt.Sprintf("%s/healthz", c.baseURL)
	resp, err := c.client.Get(url)
	if err != nil {
		return false, fmt.Errorf("health check failed to connect to %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d from %s", resp.StatusCode, url)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return false, fmt.Errorf("failed to decode response from %s: %w", url, err)
	}

	return response.Success, nil
}
