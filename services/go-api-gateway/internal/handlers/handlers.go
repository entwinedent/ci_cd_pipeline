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

// RegisterRoutes registers all HTTP routes
func (h *HTTPHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/healthz", h.HealthCheck).Methods("GET")
	router.HandleFunc("/robots.txt", h.RobotsTxt).Methods("GET")
	router.HandleFunc("/sitemap.xml", h.Sitemap).Methods("GET")
	router.HandleFunc("/api/v1/data/{key}", h.GetData).Methods("GET")
	router.HandleFunc("/api/v1/data/{key}", h.SetData).Methods("POST")
	router.HandleFunc("/api/v1/data/{key}", h.DeleteData).Methods("DELETE")
}

// HealthCheck returns the health status of the service
func (h *HTTPHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	healthy, err := h.dataStore.HealthCheck()
	if err != nil {
		log.Printf("Health check failed: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]bool{"healthy": false})
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"healthy": healthy})
}

// RobotsTxt returns the robots.txt content
func (h *HTTPHandler) RobotsTxt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User-agent: *\nDisallow: /api/\n"))
}

// Sitemap returns the sitemap.xml content
func (h *HTTPHandler) Sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>http://localhost:8080/healthz</loc>
    <lastmod>` + time.Now().Format("2006-01-02") + `</lastmod>
    <changefreq>daily</changefreq>
    <priority>0.8</priority>
  </url>
</urlset>`))
}

// LivenessCheck returns the liveness status of the gateway (always healthy if running)
func (h *HTTPHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// ReadinessCheck returns the readiness status of the gateway
func (h *HTTPHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check data store health
	healthy, err := h.dataStore.HealthCheck()
	if err != nil || !healthy {
		w.WriteHeader(http.StatusServiceUnavailable)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"ready": false,
			"error": err.Error(),
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"ready":     true,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// SetData handles POST /api/v1/data/{key} requests
func (h *HTTPHandler) SetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key is required",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	var req struct {
		Value      []byte `json:"value"`
		TTLSeconds int64  `json:"ttl_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	success, err := h.dataStore.Set(key, req.Value, req.TTLSeconds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": "Data stored successfully",
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// GetData handles GET /api/v1/data/{key} requests
func (h *HTTPHandler) GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key is required",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	value, found, err := h.dataStore.Get(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key not found",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"value":   value,
		"found":   true,
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// DeleteData handles DELETE /api/v1/data/{key} requests
func (h *HTTPHandler) DeleteData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "key is required",
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	success, err := h.dataStore.Delete(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": "Data deleted successfully",
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}
