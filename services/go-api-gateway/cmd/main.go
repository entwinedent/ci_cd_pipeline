package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/config"
	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/grpc"
	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/handlers"
	"github.com/username/ci-cd-pipeline/go-api-gateway/internal/middleware"
)

type Server struct {
	config     *config.Config
	router     *mux.Router
	httpServer *http.Server
	dataStore  handlers.DataStoreClient
}

func main() {
	if err := Run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func Run() error {
	cfg := config.LoadConfig()

	// Initialize gRPC client
	dataStoreClient, err := grpc.NewDataStoreClient(cfg.DataStoreTarget)
	if err != nil {
		return fmt.Errorf("failed to connect to data store: %w", err)
	}
	defer dataStoreClient.Close()

	log.Printf("Connected to data store at %s", cfg.DataStoreTarget)

	server := NewServer(cfg, dataStoreClient)

	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("Starting API Gateway on port %s", cfg.Port)
		if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("server failed to start: %w", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Shutting down server...")
	case err := <-serverErr:
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

func NewServer(cfg *config.Config, dataStoreClient handlers.DataStoreClient) *Server {
	router := mux.NewRouter()

	server := &Server{
		config:    cfg,
		router:    router,
		dataStore: dataStoreClient,
	}

	server.setupRoutes()

	server.httpServer = &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	return server
}

func (s *Server) setupRoutes() {
	// Apply middleware
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.CORS)
	s.router.Use(middleware.Recovery)
	s.router.Use(middleware.RequestID)

	// Create handler
	handler := handlers.NewHTTPHandler(s.dataStore)

	// Health check endpoints
	s.router.HandleFunc("/healthz", handler.HealthCheck).Methods("GET")
	s.router.HandleFunc("/readyz", handler.ReadinessCheck).Methods("GET")
	s.router.HandleFunc("/livez", handler.LivenessCheck).Methods("GET")

	// SEO endpoints
	s.router.HandleFunc("/robots.txt", handler.RobotsTxt).Methods("GET")
	s.router.HandleFunc("/sitemap.xml", handler.Sitemap).Methods("GET")

	// API routes
	s.router.HandleFunc("/api/v1/data/{key}", handler.GetData).Methods("GET")
	s.router.HandleFunc("/api/v1/data/{key}", handler.SetData).Methods("POST")
	s.router.HandleFunc("/api/v1/data/{key}", handler.DeleteData).Methods("DELETE")
}
