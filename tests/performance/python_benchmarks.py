"""
Performance benchmarks for Python Telemetry Collector
"""

import pytest
import time
import statistics
from typing import List, Dict, Any


class BenchmarkResult:
    """Container for benchmark results"""
    def __init__(self, name: str, durations: List[float]):
        self.name = name
        self.durations = durations
    
    @property
    def mean(self) -> float:
        return statistics.mean(self.durations)
    
    @property
    def median(self) -> float:
        return statistics.median(self.durations)
    
    @property
    def p95(self) -> float:
        sorted_durations = sorted(self.durations)
        index = int(len(sorted_durations) * 0.95)
        return sorted_durations[index]
    
    @property
    def p99(self) -> float:
        sorted_durations = sorted(self.durations)
        index = int(len(sorted_durations) * 0.99)
        return sorted_durations[index]
    
    def __str__(self) -> str:
        return (f"{self.name}:\n"
                f"  Mean: {self.mean:.4f}s\n"
                f"  Median: {self.median:.4f}s\n"
                f"  P95: {self.p95:.4f}s\n"
                f"  P99: {self.p99:.4f}s")


def benchmark_log_ingestion(iterations: int = 1000) -> BenchmarkResult:
    """Benchmark log ingestion performance"""
    from telemetry_client import TelemetryClient
    
    client = TelemetryClient("http://localhost:8000")
    durations = []
    
    for i in range(iterations):
        log_data = {
            "service": "benchmark-service",
            "level": "info",
            "message": f"Benchmark log {i}",
            "timestamp": time.time(),
            "metadata": {
                "request_id": f"req-{i}",
                "duration_ms": i % 100
            }
        }
        
        start = time.time()
        client.ingest_log(log_data)
        duration = time.time() - start
        durations.append(duration)
    
    return BenchmarkResult("Log Ingestion", durations)


def benchmark_metrics_query(iterations: int = 1000) -> BenchmarkResult:
    """Benchmark metrics query performance"""
    from telemetry_client import TelemetryClient
    
    client = TelemetryClient("http://localhost:8000")
    durations = []
    
    for i in range(iterations):
        start = time.time()
        client.query_metrics("benchmark-service", "1h")
        duration = time.time() - start
        durations.append(duration)
    
    return BenchmarkResult("Metrics Query", durations)


def benchmark_anomaly_detection(iterations: int = 100) -> BenchmarkResult:
    """Benchmark anomaly detection performance"""
    from telemetry_client import TelemetryClient
    
    client = TelemetryClient("http://localhost:8000")
    durations = []
    
    # Pre-populate with metrics
    for i in range(1000):
        client.ingest_metric("benchmark-service", "test_metric", i % 100)
    
    for i in range(iterations):
        start = time.time()
        client.detect_anomalies("benchmark-service")
        duration = time.time() - start
        durations.append(duration)
    
    return BenchmarkResult("Anomaly Detection", durations)


def benchmark_concurrent_operations(iterations: int = 1000) -> BenchmarkResult:
    """Benchmark concurrent operations"""
    import concurrent.futures
    from telemetry_client import TelemetryClient
    
    client = TelemetryClient("http://localhost:8000")
    durations = []
    
    def ingest_log(i: int) -> float:
        log_data = {
            "service": "benchmark-service",
            "level": "info",
            "message": f"Concurrent log {i}",
            "timestamp": time.time()
        }
        start = time.time()
        client.ingest_log(log_data)
        return time.time() - start
    
    with concurrent.futures.ThreadPoolExecutor(max_workers=10) as executor:
        futures = [executor.submit(ingest_log, i) for i in range(iterations)]
        for future in concurrent.futures.as_completed(futures):
            durations.append(future.result())
    
    return BenchmarkResult("Concurrent Operations", durations)


def test_log_ingestion_benchmark():
    """Run log ingestion benchmark"""
    result = benchmark_log_ingestion(100)
    print(result)
    assert result.mean < 0.1, f"Mean duration {result.mean}s exceeds threshold"


def test_metrics_query_benchmark():
    """Run metrics query benchmark"""
    result = benchmark_metrics_query(100)
    print(result)
    assert result.mean < 0.05, f"Mean duration {result.mean}s exceeds threshold"


def test_anomaly_detection_benchmark():
    """Run anomaly detection benchmark"""
    result = benchmark_anomaly_detection(50)
    print(result)
    assert result.mean < 0.5, f"Mean duration {result.mean}s exceeds threshold"


def test_concurrent_operations_benchmark():
    """Run concurrent operations benchmark"""
    result = benchmark_concurrent_operations(100)
    print(result)
    assert result.mean < 0.2, f"Mean duration {result.mean}s exceeds threshold"


def test_all_benchmarks():
    """Run all benchmarks and generate report"""
    print("=" * 50)
    print("Performance Benchmark Report")
    print("=" * 50)
    
    results = [
        benchmark_log_ingestion(100),
        benchmark_metrics_query(100),
        benchmark_anomaly_detection(50),
        benchmark_concurrent_operations(100),
    ]
    
    print("\n" + "=" * 50)
    print("Summary")
    print("=" * 50)
    
    for result in results:
        print(result)


if __name__ == "__main__":
    test_all_benchmarks()
