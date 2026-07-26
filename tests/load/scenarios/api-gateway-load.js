import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

export const options = {
  stages: [
    { duration: '30s', target: 10 },   // Ramp up to 10 users
    { duration: '1m', target: 50 },    // Ramp up to 50 users
    { duration: '2m', target: 100 },   // Ramp up to 100 users
    { duration: '2m', target: 100 },   // Stay at 100 users
    { duration: '1m', target: 50 },    // Ramp down to 50 users
    { duration: '30s', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'], // 95% of requests < 500ms, 99% < 1s
    http_req_failed: ['rate<0.01'],                   // Error rate < 1%
    errors: ['rate<0.01'],
  },
};

const BASE_URL = __ENV.GATEWAY_URL || 'http://localhost:8080';

export default function () {
  // Test health endpoint
  let healthRes = http.get(`${BASE_URL}/healthz`);
  check(healthRes, {
    'health status is 200': (r) => r.status === 200,
    'health response time < 100ms': (r) => r.timings.duration < 100,
  }) || errorRate.add(1);

  // Test set operation
  const key = `test_key_${__VU}_${__ITER}`;
  let setRes = http.post(
    `${BASE_URL}/api/v1/data/${key}`,
    JSON.stringify({ value: `test_value_${__ITER}` }),
    {
      headers: { 'Content-Type': 'application/json' },
    }
  );
  check(setRes, {
    'set status is 200': (r) => r.status === 200,
    'set response time < 200ms': (r) => r.timings.duration < 200,
    'set has success': (r) => r.json('success') === true,
  }) || errorRate.add(1);

  // Test get operation
  let getRes = http.get(`${BASE_URL}/api/v1/data/${key}`);
  check(getRes, {
    'get status is 200': (r) => r.status === 200,
    'get response time < 200ms': (r) => r.timings.duration < 200,
    'get has value': (r) => r.json('value') === `test_value_${__ITER}`,
  }) || errorRate.add(1);

  // Test delete operation
  let deleteRes = http.del(`${BASE_URL}/api/v1/data/${key}`);
  check(deleteRes, {
    'delete status is 200': (r) => r.status === 200,
    'delete response time < 200ms': (r) => r.timings.duration < 200,
    'delete has success': (r) => r.json('success') === true,
  }) || errorRate.add(1);

  sleep(1); // Pause between iterations
}

export function handleSummary(data) {
  return {
    'summary.json': JSON.stringify(data),
    stdout: textSummary(data, { indent: ' ', enableColors: true }),
  };
}
