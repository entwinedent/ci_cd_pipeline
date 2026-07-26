import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '30s', target: 5 },    // Ramp up to 5 users
    { duration: '1m', target: 20 },    // Ramp up to 20 users
    { duration: '2m', target: 50 },    // Ramp up to 50 users
    { duration: '2m', target: 50 },    // Stay at 50 users
    { duration: '1m', target: 20 },    // Ramp down to 20 users
    { duration: '30s', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<300', 'p(99)<500'],
    http_req_failed: ['rate<0.02'],
    errors: ['rate<0.02'],
  },
};

const BASE_URL = __ENV.TELEMETRY_URL || 'http://localhost:8000';

export default function () {
  // Test health endpoint
  let healthRes = http.get(`${BASE_URL}/healthz`);
  check(healthRes, {
    'health status is 200': (r) => r.status === 200,
    'health response time < 100ms': (r) => r.timings.duration < 100,
  }) || errorRate.add(1);

  // Test log ingestion
  const logData = {
    service: `test-service-${__VU}`,
    level: 'info',
    message: `Test log message ${__ITER}`,
    timestamp: new Date().toISOString(),
    metadata: {
      request_id: `req-${__VU}-${__ITER}`,
      duration_ms: Math.random() * 100,
    },
  };
  let logRes = http.post(
    `${BASE_URL}/api/v1/logs`,
    JSON.stringify(logData),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
  check(logRes, {
    'log status is 200': (r) => r.status === 200,
    'log response time < 300ms': (r) => r.timings.duration < 300,
    'log has success': (r) => r.json('success') === true,
  }) || errorRate.add(1);

  // Test metrics query
  let metricsRes = http.get(
    `${BASE_URL}/api/v1/metrics?service=test-service-${__VU}&time_range=1h`
  );
  check(metricsRes, {
    'metrics status is 200': (r) => r.status === 200,
    'metrics response time < 200ms': (r) => r.timings.duration < 200,
    'metrics has data': (r) => 'metrics' in r.json(),
  }) || errorRate.add(1);

  sleep(1);
}

export function handleSummary(data) {
  return {
    'telemetry-summary.json': JSON.stringify(data),
  };
}
