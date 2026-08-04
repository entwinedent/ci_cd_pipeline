package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Message: "Service is healthy",
		})
	})

	http.HandleFunc("/api/v1/data/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		switch r.Method {
		case http.MethodPut:
			json.NewEncoder(w).Encode(Response{
				Success: true,
				Message: "Data stored successfully",
			})
		case http.MethodGet:
			json.NewEncoder(w).Encode(Response{
				Success: true,
				Message: "Data retrieved successfully",
				Value:   "test-value",
			})
		case http.MethodDelete:
			json.NewEncoder(w).Encode(Response{
				Success: true,
				Message: "Data deleted successfully",
			})
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	port := "50051"
	fmt.Printf("Mock data store server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
