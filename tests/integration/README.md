# Integration Tests

This directory contains integration tests for the CI/CD pipeline services, testing end-to-end communication between the Go API Gateway, Rust Data Store, and Python Telemetry Collector.

## Test Coverage

### Go API Gateway Tests
- Health check endpoint
- Readiness check endpoint
- Set data operation
- Get data operation
- Delete data operation

### Rust Data Store Tests (gRPC)
- Health check via gRPC
- Set operation via gRPC
- Get operation via gRPC
- Delete operation via gRPC

### Python Telemetry Tests
- Health check endpoint
- Readiness check endpoint
- Log ingestion endpoint
- Metrics query endpoint
- Anomaly detection endpoint

### End-to-End Tests
- Complete data flow through all services
- Service communication validation
- Data consistency checks

## Prerequisites

### Services Running

All services must be running before executing integration tests:

```bash
# Option 1: Docker Compose
docker-compose up -d

# Option 2: Kind Cluster
make kind-setup
make deploy
```

### Environment Variables

Optional environment variables to customize service URLs:

```bash
export GATEWAY_URL="http://localhost:8080"
export TELEMETRY_URL="http://localhost:8000"
export DATA_STORE_URL="localhost:50051"
```

## Running Tests

### Local Development

```bash
# Install dependencies
pip install -r requirements.txt

# Run all integration tests
pytest test_integration.py -v

# Run specific test class
pytest test_integration.py::TestGoAPIGateway -v

# Run specific test
pytest test_integration.py::TestGoAPIGateway::test_health_check -v

# Run with coverage
pytest test_integration.py --cov=. --cov-report=html
```

### CI/CD Pipeline

Integration tests are automatically run in the CI pipeline:

```bash
make test-integration
```

## Test Fixtures

### wait_for_services

Automatically waits for all services to be healthy before running tests:
- Checks API Gateway health endpoint
- Checks Telemetry health endpoint
- Checks Data Store gRPC health check
- Retries up to 30 times with 2-second delays

### cleanup_data

Automatically cleans up test data after each test to prevent test pollution.

## Troubleshooting

### Services Not Ready

If tests fail due to services not being ready:

```bash
# Check service status
make status

# Check service logs
make logs

# Restart services
docker-compose restart
# or
kubectl rollout restart deployment/go-api-gateway
kubectl rollout restart deployment/rust-data-store
kubectl rollout restart deployment/python-telemetry
```

### Connection Refused

If tests fail with connection refused errors:

```bash
# Verify services are running
docker ps
# or
kubectl get pods

# Check port forwarding
make port-forward
```

### gRPC Errors

If gRPC tests fail:

```bash
# Verify Data Store is running
kubectl get pods -l app=rust-data-store

# Check Data Store logs
kubectl logs -l app=rust-data-store

# Test gRPC connection manually
grpcurl -plaintext localhost:50051 list
```

## Adding New Tests

When adding new integration tests:

1. Add test methods to appropriate test class
2. Follow naming convention: `test_<feature>`
3. Include descriptive docstrings
4. Use fixtures for common setup/teardown
5. Update this README with new test coverage

## Continuous Integration

Integration tests run in GitHub Actions workflow:

```yaml
- name: Run Integration Tests
  run: |
    docker-compose up -d
    make test-integration
    docker-compose down
```

## Best Practices

- Use descriptive test names
- Keep tests independent
- Use fixtures for common setup
- Clean up test data after tests
- Test both success and failure cases
- Include assertions for all expected behaviors
- Add comments for complex test logic
