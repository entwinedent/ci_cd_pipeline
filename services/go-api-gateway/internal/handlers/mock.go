package handlers

import "fmt"

// MockDataStoreClient is a mock implementation for testing
type MockDataStoreClient struct {
	healthy          bool
	getNotFound      bool
	setError         bool
	getError         bool
	deleteError      bool
	healthCheckError bool
}

// NewMockDataStoreClient creates a new mock client
func NewMockDataStoreClient(healthy bool) *MockDataStoreClient {
	return &MockDataStoreClient{healthy: healthy}
}

func (m *MockDataStoreClient) Set(key string, value []byte, ttl int64) (bool, error) {
	if m.setError {
		return false, fmt.Errorf("set error")
	}
	return true, nil
}

func (m *MockDataStoreClient) Get(key string) ([]byte, bool, error) {
	if m.getError {
		return nil, false, fmt.Errorf("get error")
	}
	if m.getNotFound {
		return nil, false, nil
	}
	return []byte("test"), true, nil
}

func (m *MockDataStoreClient) Delete(key string) (bool, error) {
	if m.deleteError {
		return false, fmt.Errorf("delete error")
	}
	return true, nil
}

func (m *MockDataStoreClient) HealthCheck() (bool, error) {
	if m.healthCheckError {
		return false, fmt.Errorf("health check error")
	}
	if !m.healthy {
		return false, fmt.Errorf("data store is unhealthy")
	}
	return m.healthy, nil
}
