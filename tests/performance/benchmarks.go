package performance

import (
	"testing"
	"time"
)

// BenchmarkGoAPISet benchmarks the Go API Gateway set operation
func BenchmarkGoAPISet(b *testing.B) {
	// Setup benchmark client
	client := setupBenchmarkClient()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		value := fmt.Sprintf("bench_value_%d", i)
		_, err := client.SetData(key, value)
		if err != nil {
			b.Fatalf("SetData failed: %v", err)
		}
	}
}

// BenchmarkGoAPIGet benchmarks the Go API Gateway get operation
func BenchmarkGoAPIGet(b *testing.B) {
	client := setupBenchmarkClient()
	
	// Pre-populate data
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("bench_key_%d", i)
		client.SetData(key, fmt.Sprintf("bench_value_%d", i))
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("bench_key_%d", i%1000)
		_, err := client.GetData(key)
		if err != nil {
			b.Fatalf("GetData failed: %v", err)
		}
	}
}

// BenchmarkGoAPIConcurrent benchmarks concurrent API operations
func BenchmarkGoAPIConcurrent(b *testing.B) {
	client := setupBenchmarkClient()
	
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("concurrent_key_%d", i)
			client.SetData(key, fmt.Sprintf("concurrent_value_%d", i))
			i++
		}
	})
}

// BenchmarkRustGRPCSet benchmarks Rust Data Store gRPC set operation
func BenchmarkRustGRPCSet(b *testing.B) {
	conn := setupGRPCConnection()
	defer conn.Close()
	client := datastore.NewDataStoreServiceClient(conn)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		
		_, err := client.Set(ctx, &datastore.SetRequest{
			Key:   []byte(fmt.Sprintf("bench_key_%d", i)),
			Value: []byte(fmt.Sprintf("bench_value_%d", i)),
		})
		if err != nil {
			b.Fatalf("Set failed: %v", err)
		}
	}
}

// BenchmarkRustGRPCGet benchmarks Rust Data Store gRPC get operation
func BenchmarkRustGRPCGet(b *testing.B) {
	conn := setupGRPCConnection()
	defer conn.Close()
	client := datastore.NewDataStoreServiceClient(conn)
	
	// Pre-populate data
	for i := 0; i < 1000; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		client.Set(ctx, &datastore.SetRequest{
			Key:   []byte(fmt.Sprintf("bench_key_%d", i)),
			Value: []byte(fmt.Sprintf("bench_value_%d", i)),
		})
		cancel()
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		
		key := fmt.Sprintf("bench_key_%d", i%1000)
		_, err := client.Get(ctx, &datastore.GetRequest{
			Key: []byte(key),
		})
		if err != nil {
			b.Fatalf("Get failed: %v", err)
		}
	}
}

// BenchmarkPythonTelemetryLogIngestion benchmarks Python telemetry log ingestion
func BenchmarkPythonTelemetryLogIngestion(b *testing.B) {
	client := setupTelemetryClient()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logData := map[string]interface{}{
			"service":  "benchmark-service",
			"level":    "info",
			"message":  fmt.Sprintf("Benchmark log %d", i),
			"timestamp": time.Now().Format(time.RFC3339),
		}
		_, err := client.IngestLog(logData)
		if err != nil {
			b.Fatalf("IngestLog failed: %v", err)
		}
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(b *testing.B) {
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		data := make([]byte, 1024)
		_ = data
	}
}

// Helper functions
func setupBenchmarkClient() *APIClient {
	// Setup benchmark API client
	return &APIClient{
		BaseURL: "http://localhost:8080",
	}
}

func setupGRPCConnection() *grpc.ClientConn {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}
	return conn
}

func setupTelemetryClient() *TelemetryClient {
	return &TelemetryClient{
		BaseURL: "http://localhost:8000",
	}
}
