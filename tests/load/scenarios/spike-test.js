import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');
const API_GATEWAY_URL = __ENV.API_GATEWAY_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '10s', target: 10 },   // Normal load
    { duration: '10s', target: 500 },  // Massive spike
    { duration: '20s', target: 500 },  // Sustained spike
    { duration: '10s', target: 10 },   // Recovery
  ],
  thresholds: {
    http_req_duration: ['p(95)<1000'], // More lenient during spike
    errors: ['rate<0.1'], // Allow higher error rate during spike
  },
};

export default function () {
  const healthRes = http.get(`${API_GATEWAY_URL}/healthz`);
  check(healthRes, {
    'health check status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  const testKey = `spike-test-${__VU}-${Date.now()}`;
  const setRes = http.post(`${API_GATEWAY_URL}/api/v1/data`, JSON.stringify({
    key: testKey,
    value: 'spike-test-value',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(setRes, {
    'set data status 201': (r) => r.status === 201,
  }) || errorRate.add(1);

  sleep(0.1); // Faster requests during spike
}
