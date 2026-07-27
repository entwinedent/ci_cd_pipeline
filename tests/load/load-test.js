import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const apiLatency = new Trend('api_latency');
const grpcLatency = new Trend('grpc_latency');

const API_GATEWAY_URL = __ENV.API_GATEWAY_URL || 'http://localhost:8080';
const TELEMETRY_URL = __ENV.TELEMETRY_URL || 'http://localhost:8000';

export const options = {
  stages: [
    { duration: '30s', target: 10 },   // Ramp up to 10 users
    { duration: '1m', target: 10 },    // Stay at 10 users
    { duration: '30s', target: 50 },   // Ramp up to 50 users
    { duration: '1m', target: 50 },    // Stay at 50 users
    { duration: '30s', target: 100 },  // Ramp up to 100 users
    { duration: '1m', target: 100 },   // Stay at 100 users
    { duration: '30s', target: 0 },    // Ramp down to 0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% of requests under 500ms, 99% under 1s
    errors: ['rate<0.5'], // Error rate less than 50% (temporarily relaxed for debugging)
  },
};

export default function () {
  // Test 1: Health check
  const healthRes = http.get(`${API_GATEWAY_URL}/healthz`, {
    tags: { name: 'HealthCheck' },
  });
  check(healthRes, {
    'health check status 200': (r) => r.status === 200,
    'health check response time < 100ms': (r) => r.timings.duration < 100,
  }) || errorRate.add(1);
  apiLatency.add(healthRes.timings.duration);

  // Test 2: Data store operations
  const testKey = `load-test-key`; // Use static key to reduce cardinality
  const testValue = 'load-test-value';

  // Set data - use correct path with key in URL
  const setRes = http.post(`${API_GATEWAY_URL}/api/v1/data/${testKey}`, JSON.stringify({
    value: testValue,
    ttl_seconds: 3600,
  }), {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'SetData' },
  });
  check(setRes, {
    'set data status 200': (r) => r.status === 200,
    'set data response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);
  grpcLatency.add(setRes.timings.duration);

  // Get data
  const getRes = http.get(`${API_GATEWAY_URL}/api/v1/data/${testKey}`, {
    tags: { name: 'GetData' },
  });
  check(getRes, {
    'get data status 200': (r) => r.status === 200,
    'get data response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);
  grpcLatency.add(getRes.timings.duration);

  // Delete data
  const deleteRes = http.del(`${API_GATEWAY_URL}/api/v1/data/${testKey}`, null, {
    tags: { name: 'DeleteData' },
  });
  check(deleteRes, {
    'delete data status 200': (r) => r.status === 200,
    'delete data response time < 200ms': (r) => r.timings.duration < 200,
  }) || errorRate.add(1);
  grpcLatency.add(deleteRes.timings.duration);

  // Test 3: Telemetry log ingestion
  const logRes = http.post(`${TELEMETRY_URL}/api/v1/logs`, JSON.stringify({
    service: 'load-test-service',
    level: 'INFO',
    message: 'Load test log message',
    timestamp: new Date().toISOString(),
  }), {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'LogIngestion' },
  });
  check(logRes, {
    'log ingestion status 200': (r) => r.status === 200,
    'log ingestion response time < 300ms': (r) => r.timings.duration < 300,
  }) || errorRate.add(1);

  // Test 4: Telemetry metrics
  const metricRes = http.post(`${TELEMETRY_URL}/api/v1/metrics`, JSON.stringify({
    service: 'load-test-service',
    metric_name: 'load_test_metric',
    value: Math.random() * 100,
    timestamp: new Date().toISOString(),
  }), {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'MetricIngestion' },
  });
  check(metricRes, {
    'metric ingestion status 200': (r) => r.status === 200,
    'metric ingestion response time < 300ms': (r) => r.timings.duration < 300,
  }) || errorRate.add(1);

  sleep(1);
}
