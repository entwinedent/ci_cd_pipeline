package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// DataStoreClient defines the interface for the gRPC data store client
type DataStoreClient interface {
	Set(key string, value []byte, ttl int64) (bool, error)
	Get(key string) ([]byte, bool, error)
	Delete(key string) (bool, error)
	HealthCheck() (bool, error)
}

// HTTPHandler handles HTTP requests for the API Gateway
type HTTPHandler struct {
	dataStore DataStoreClient
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(dataStore DataStoreClient) *HTTPHandler {
	return &HTTPHandler{
		dataStore: dataStore,
	}
}

// HealthCheck returns the health status of the gateway
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check data store health
	healthy, err := h.dataStore.HealthCheck()
	if err != nil {
		log.Printf("Health check error: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	if !healthy {
		log.Printf("Health check returned unhealthy")
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "unhealthy",
			"error":  "data store reported unhealthy",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	log.Printf("Health check successful")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// LivenessCheck returns the liveness status of the gateway (always healthy if running)
func (h *HTTPHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// ReadinessCheck returns the readiness status of the gateway
func (h *HTTPHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check data store health
	healthy, err := h.dataStore.HealthCheck()
	if err != nil || !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ready": false,
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ready":     true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// SetData handles POST /api/v1/data/{key} requests
func (h *HTTPHandler) SetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key is required",
		})
		return
	}

	var req struct {
		Value      []byte `json:"value"`
		TTLSeconds int64  `json:"ttl_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	success, err := h.dataStore.Set(key, req.Value, req.TTLSeconds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": "Data stored successfully",
	})
}

// GetData handles GET /api/v1/data/{key} requests
func (h *HTTPHandler) GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key is required",
		})
		return
	}

	value, found, err := h.dataStore.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"value":   value,
		"found":   true,
	})
}

// DeleteData handles DELETE /api/v1/data/{key} requests
func (h *HTTPHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key is required",
		})
		return
	}

	success, err := h.dataStore.Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": "Data deleted successfully",
	})
}
