# Go API Gateway

High-throughput HTTP/gRPC service serving as the external entrypoint for the CI/CD pipeline microservices architecture.

## Architecture

The API Gateway provides:
- HTTP REST API for external clients
- gRPC client for communicating with Rust Data Store
- Feature flag integration (Unleash/LaunchDarkly)
- Request routing and load balancing
- Rate limiting and circuit breaking

## Routing Logic

### HTTP Endpoints

- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe
- `GET /api/v1/data/{key}` - Retrieve data from data store
- `POST /api/v1/data/{key}` - Store data in data store
- `DELETE /api/v1/data/{key}` - Delete data from data store

### gRPC Communication

The gateway communicates with the Rust Data Store via gRPC on port 50051:
```protobuf
service DataStoreService {
  rpc Set(SetRequest) returns (SetResponse);
  rpc Get(GetRequest) returns (GetResponse);
  rpc Delete(DeleteRequest) returns (DeleteResponse);
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}
```

## Feature Flag Integration

The gateway uses a provider-agnostic feature flag abstraction:

### Supported Providers
- **Unleash** (self-hosted)
- **LaunchDarkly** (commercial)

### Usage Example

```go
cfg := featureflags.Config{
    Provider: featureflags.ProviderUnleash,
}
provider, _ := featureflags.NewProvider(cfg)
ff := featureflags.NewFeatureFlags(provider)

if ff.NewAPIEndpointEnabled(ctx) {
    // Use new endpoint implementation
}
```

### Environment Variables
- `FEATURE_FLAG_PROVIDER` - Select provider (`unleash` or `launchdarkly`)
- `UNLEASH_URL` - Unleash server URL
- `UNLEASH_API_TOKEN` - Unleash API token
- `LAUNCHDARKLY_SDK_KEY` - LaunchDarkly SDK key

## Rate Limiting

Rate limiting is implemented using token bucket algorithm:
- Default: 100 requests per minute per IP
- Configurable via environment variables
- Circuit breaking on service degradation

## Configuration

### Environment Variables
- `PORT` - Service listening port (default: 8080)
- `DATA_STORE_TARGET` - gRPC endpoint for Rust service (default: `rust-data-store:50051`)
- `LOG_LEVEL` - Logging level (default: `info`)
- `RATE_LIMIT_RPM` - Rate limit requests per minute (default: 100)

### Service Specifications
- **Port**: 8080
- **Health**: `/healthz` (liveness), `/readyz` (readiness)
- **Metrics**: Prometheus metrics on `/metrics`

## Development

### Build
```bash
cd services/go-api-gateway
go build -o bin/api-gateway ./cmd
```

### Run
```bash
./bin/api-gateway
```

### Test
```bash
go test ./...
go test -race ./...
go test -cover ./...
```

### Docker Build
```bash
docker build -t go-api-gateway:latest .
```

## Performance

- **Image Size**: < 15MB
- **Response Time**: < 100ms (p95)
- **Throughput**: > 10,000 RPS
- **Memory**: < 50MB

## Dependencies

- `github.com/gorilla/mux` - HTTP router
- `google.golang.org/grpc` - gRPC framework
- `github.com/spiffe/go-spiffe/v2` - SPIFFE/SPIRE integration

## Security

- SPIFFE/SPIRE mTLS for service-to-service communication
- Input validation and sanitization
- Rate limiting to prevent abuse
- Structured logging with sensitive data redaction
