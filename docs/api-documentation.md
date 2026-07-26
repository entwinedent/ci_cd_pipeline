# API Documentation

## Overview

This CI/CD pipeline implements a polyglot microservices architecture with three core services communicating via gRPC and HTTP protocols. The system is designed for high-performance, zero-trust communication with comprehensive observability.

## Services

### 1. Go API Gateway

**Port:** 8080  
**Protocol:** HTTP/1.1 (external), gRPC (internal)  
**Purpose:** External API entrypoint, request validation, and routing

#### Endpoints

**Health Checks**
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe

**Data Operations**
- `GET /api/v1/data/{key}` - Retrieve value by key
- `POST /api/v1/data/{key}` - Store value with JSON body
- `DELETE /api/v1/data/{key}` - Delete value by key

**Request/Response Format**

```json
// POST /api/v1/data/{key}
{
  "value": "string",
  "ttl": 3600  // Optional time-to-live in seconds
}

// Response
{
  "success": true,
  "message": "Data stored successfully"
}
```

#### gRPC Client Configuration

The Go API Gateway acts as a gRPC client to the Rust Data Store:

```go
// Connection configuration
conn, err := grpc.Dial("rust-data-store:50051", 
    grpc.WithTransportCredentials(insecure.NewCredentials()))
client := datastore.NewDataStoreServiceClient(conn)
```

### 2. Rust Data Store

**Port:** 50051  
**Protocol:** gRPC  
**Purpose:** High-performance in-memory storage with thread-safe operations

#### gRPC Service Definition

```protobuf
service DataStoreService {
  rpc Set(SetRequest) returns (SetResponse);
  rpc Get(GetRequest) returns (GetResponse);
  rpc Delete(DeleteRequest) returns (DeleteResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}
```

#### Message Types

**SetRequest**
```protobuf
message SetRequest {
  string key = 1;
  string value = 2;
  int32 ttl = 3;  // Time-to-live in seconds
}
```

**GetRequest**
```protobuf
message GetRequest {
  string key = 1;
}
```

**GetResponse**
```protobuf
message GetResponse {
  bool success = 1;
  string value = 2;
  string error = 3;
}
```

#### Performance Characteristics

- **Concurrent Access**: Thread-safe using DashMap
- **Memory Efficiency**: Lock-free concurrent hash map
- **TTL Management**: Binary min-heap for expiration
- **Latency**: < 10ms for read operations

### 3. Python Telemetry

**Port:** 8000  
**Protocol:** HTTP/1.1  
**Purpose:** Log collection, metrics aggregation, and AI anomaly detection

#### Endpoints

**Health Checks**
- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe

**Log Ingestion**
- `POST /api/v1/logs` - Ingest structured logs
- `GET /api/v1/logs` - Query logs with filters

**Metrics**
- `GET /api/v1/metrics` - Retrieve aggregated metrics
- `POST /api/v1/metrics` - Push custom metrics

**Anomaly Detection**
- `GET /api/v1/anomalies` - Get detected anomalies
- `POST /api/v1/anomalies/analyze` - Trigger analysis

#### Log Format

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "INFO",
  "service": "go-api-gateway",
  "message": "Request processed successfully",
  "metadata": {
    "request_id": "abc123",
    "duration_ms": 45
  }
}
```

#### AI Anomaly Detection

The Python service uses machine learning to detect anomalies:

```python
# Anomaly detection pipeline
def detect_anomalies(metrics):
    # Time-series analysis
    # Statistical outlier detection
    # Pattern recognition
    return anomalies
```

## Service Communication Flow

```
External Client
       │
       │ HTTP Request
       ▼
┌──────────────────┐
│  Go API Gateway  │
│  Port 8080       │
└────────┬─────────┘
         │
         │ gRPC Request
         ▼
┌──────────────────┐
│  Rust Data Store │
│  Port 50051      │
└────────┬─────────┘
         │
         │ gRPC Response
         ▼
┌──────────────────┐
│  Go API Gateway  │
└────────┬─────────┘
         │
         │ HTTP Response
         ▼
   External Client

Parallel Telemetry Flow:
┌──────────────────┐
│  Go API Gateway  │
└────────┬─────────┘
         │
         │ Log/Metric Data
         ▼
┌──────────────────┐
│ Python Telemetry │
│  Port 8000       │
└──────────────────┘
```

## Authentication & Security

### Zero-Trust Architecture

- **SPIFFE/SPIRE**: Cryptographic workload identity
- **mTLS**: Mutual TLS for service-to-service communication
- **SVID Certificates**: Short-lived X.509 certificates

### Certificate Rotation

- **Validity**: 1 hour default
- **Rotation**: Automatic via SPIRE agent
- **Revocation**: Immediate on workload termination

## Error Handling

### Standard Error Response Format

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameter",
    "details": "Key must be alphanumeric"
  }
}
```

### Error Codes

- `VALIDATION_ERROR` - Invalid request parameters
- `NOT_FOUND` - Resource not found
- `INTERNAL_ERROR` - Server-side error
- `RATE_LIMIT_EXCEEDED` - Too many requests
- `AUTHENTICATION_FAILED` - Invalid credentials

## Rate Limiting

### Default Limits

- **Go API Gateway**: 1000 requests/minute
- **Rust Data Store**: 5000 operations/minute
- **Python Telemetry**: 2000 logs/second

### Rate Limit Headers

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 950
X-RateLimit-Reset: 1642245600
```

## Observability

### Distributed Tracing

- **Trace ID**: W3C Trace Context format
- **Span ID**: Per-operation tracing
- **Propagation**: HTTP headers and gRPC metadata

### Metrics

- **Counter**: Request counts, error rates
- **Gauge**: Current connections, memory usage
- **Histogram**: Request latency distribution
- **Summary**: Response time percentiles

### Logging

- **Format**: Structured JSON
- **Levels**: DEBUG, INFO, WARN, ERROR
- **Correlation**: Request IDs across services

## Deployment

### Kubernetes Configuration

**Service Discovery**
```yaml
# Go API Gateway
apiVersion: v1
kind: Service
metadata:
  name: go-api-gateway
spec:
  selector:
    app: go-api-gateway
  ports:
  - port: 8080
    targetPort: 8080
```

**Environment Variables**
```yaml
env:
  - name: DATA_STORE_TARGET
    value: "rust-data-store:50051"
  - name: LOG_LEVEL
    value: "INFO"
  - name: TELEMETRY_ENDPOINT
    value: "python-telemetry:8000"
```

## Testing

### API Testing Examples

**Go API Gateway**
```bash
# Health check
curl http://localhost:8080/healthz

# Store data
curl -X POST http://localhost:8080/api/v1/data/mykey \
  -H "Content-Type: application/json" \
  -d '{"value": "myvalue", "ttl": 3600}'

# Retrieve data
curl http://localhost:8080/api/v1/data/mykey
```

**Python Telemetry**
```bash
# Ingest logs
curl -X POST http://localhost:8000/api/v1/logs \
  -H "Content-Type: application/json" \
  -d '{
    "level": "INFO",
    "service": "test-service",
    "message": "Test log entry"
  }'

# Get anomalies
curl http://localhost:8000/api/v1/anomalies
```

## Versioning

### API Versioning Strategy

- **URL Versioning**: `/api/v1/`, `/api/v2/`
- **Backward Compatibility**: Maintain v1 while introducing v2
- **Deprecation**: 6-month notice period for deprecated endpoints

### Current Version

- **API Version**: v1.0
- **Protocol**: HTTP/1.1, gRPC 1.0
- **Compatibility**: Stable

## Performance Benchmarks

### Service Performance

| Service | Avg Latency | P95 Latency | Throughput |
|---------|-------------|-------------|------------|
| Go API Gateway | 45ms | 120ms | 1000 req/s |
| Rust Data Store | 8ms | 25ms | 5000 ops/s |
| Python Telemetry | 15ms | 50ms | 2000 logs/s |

### Resource Usage

| Service | CPU | Memory | Network |
|---------|-----|--------|---------|
| Go API Gateway | 0.5 cores | 50MB | 100MB/s |
| Rust Data Store | 0.3 cores | 25MB | 50MB/s |
| Python Telemetry | 0.8 cores | 100MB | 200MB/s |

## Troubleshooting

### Common Issues

**Connection Refused**
- Check service health endpoints
- Verify Kubernetes service discovery
- Confirm network policies allow traffic

**High Latency**
- Check resource limits
- Review network latency between services
- Monitor gRPC connection pooling

**Memory Leaks**
- Monitor memory usage over time
- Check for goroutine leaks (Go)
- Review Rust memory allocation patterns

## Documentation References

- [gRPC Documentation](https://grpc.io/docs/)
- [HTTP/1.1 Specification](https://www.rfc-editor.org/rfc/rfc7231)
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)
- [OpenTelemetry](https://opentelemetry.io/docs/)
