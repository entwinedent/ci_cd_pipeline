# Performance Benchmarks

This directory contains performance benchmarks for the CI/CD pipeline services to measure and track performance characteristics over time.

## Benchmarks

### Go Benchmarks (`benchmarks.go`)

- `BenchmarkGoAPISet` - API Gateway set operation performance
- `BenchmarkGoAPIGet` - API Gateway get operation performance  
- `BenchmarkGoAPIConcurrent` - Concurrent API operations
- `BenchmarkRustGRPCSet` - Rust Data Store gRPC set operation
- `BenchmarkRustGRPCGet` - Rust Data Store gRPC get operation
- `BenchmarkPythonTelemetryLogIngestion` - Python telemetry log ingestion
- `BenchmarkMemoryAllocation` - Memory allocation patterns

### Rust Benchmarks (`rust_benchmarks.rs`)

- `bench_set_operation` - Data Store set operation with varying payload sizes
- `bench_get_operation` - Data Store get operation with varying payload sizes
- `bench_concurrent_operations` - Concurrent read/write operations
- `bench_ttl_eviction` - TTL cleanup performance
- `bench_memory_usage` - Memory allocation benchmarks

### Python Benchmarks (`python_benchmarks.py`)

- `benchmark_log_ingestion` - Log ingestion performance
- `benchmark_metrics_query` - Metrics query performance
- `benchmark_anomaly_detection` - Anomaly detection performance
- `benchmark_concurrent_operations` - Concurrent operations

## Running Benchmarks

### Go Benchmarks

```bash
cd services/go-api-gateway
go test -bench=. -benchmem ./tests/performance/
```

### Rust Benchmarks

```bash
cd services/rust-data-store
cargo bench --bench performance
```

### Python Benchmarks

```bash
cd services/python-telemetry
python tests/performance/python_benchmarks.py
```

## Performance Targets

### Go API Gateway
- Set operation: < 1ms mean
- Get operation: < 1ms mean
- Concurrent operations: < 5ms mean
- Memory allocation: < 100KB per operation

### Rust Data Store
- Set operation: < 100μs mean
- Get operation: < 50μs mean
- TTL cleanup: < 10ms for 1000 entries
- Concurrent operations: < 200μs mean

### Python Telemetry
- Log ingestion: < 10ms mean
- Metrics query: < 5ms mean
- Anomaly detection: < 500ms mean
- Concurrent operations: < 20ms mean

## Continuous Integration

Benchmarks run in CI pipeline to detect performance regressions:

```yaml
- name: Run Performance Benchmarks
  run: |
    cd services/go-api-gateway
    go test -bench=. -benchmem ./tests/performance/
    cd ../rust-data-store
    cargo bench --bench performance
    cd ../python-telemetry
    python tests/performance/python_benchmarks.py
```

## Performance Regression Detection

Benchmarks are compared against baseline measurements:

- **Warning**: 10% degradation from baseline
- **Critical**: 25% degradation from baseline
- **Block**: 50% degradation from baseline

## Profiling

### Go Profiling

```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./tests/performance/
go tool pprof cpu.prof
go tool pprof mem.prof
```

### Rust Profiling

```bash
cargo bench --bench performance -- --profile-time 10
```

### Python Profiling

```bash
python -m cProfile -o profile.stats tests/performance/python_benchmarks.py
python -m pstats profile.stats
```

## Benchmark Results Storage

Results are stored in `tests/performance/results/`:

```
results/
├── go-benchmarks.json
├── rust-benchmarks.json
└── python-benchmarks.json
```

## Historical Tracking

Benchmark results are tracked over time to identify trends:

- Performance improvements
- Performance regressions
- Impact of code changes
- Resource utilization patterns

## Best Practices

1. **Run benchmarks consistently** - Same environment, same hardware
2. **Warm up services** - Allow services to reach steady state
3. **Multiple iterations** - Run multiple iterations for statistical significance
4. **Isolate benchmarks** - Run benchmarks independently
5. **Document changes** - Note any changes that might affect performance
6. **Monitor trends** - Track performance over time, not just single runs

## Troubleshooting

### Inconsistent Results

- Ensure services are in steady state
- Check for background processes
- Verify resource availability
- Run multiple iterations

### Slow Benchmarks

- Check service health
- Verify network connectivity
- Check resource constraints
- Profile to identify bottlenecks

### Memory Issues

- Check for memory leaks
- Verify garbage collection
- Profile memory usage
- Check resource limits
