import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const errorRate = new Rate('errors');
const memoryUsage = new Trend('memory_usage');
const API_GATEWAY_URL = __ENV.API_GATEWAY_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '2m', target: 20 },   // Ramp up
    { duration: '10m', target: 20 },  // Sustained load for 10 minutes
    { duration: '2m', target: 0 },    // Ramp down
  ],
  thresholds: {
    http_req_duration: ['p(95)<300'],
    errors: ['rate<0.02'],
  },
};

export default function () {
  const startMemory = __ENV.MEM_USAGE || 0;
  
  const healthRes = http.get(`${API_GATEWAY_URL}/healthz`);
  check(healthRes, {
    'health check status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  // Perform data operations
  const testKey = `soak-test-${__VU}-${Date.now()}`;
  const setRes = http.post(`${API_GATEWAY_URL}/api/v1/data`, JSON.stringify({
    key: testKey,
    value: 'soak-test-value',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(setRes, {
    'set data status 201': (r) => r.status === 201,
  }) || errorRate.add(1);

  const getRes = http.get(`${API_GATEWAY_URL}/api/v1/data/${testKey}`);
  check(getRes, {
    'get data status 200': (r) => r.status === 200,
  }) || errorRate.add(1);

  // Simulate memory tracking
  memoryUsage.add(Math.random() * 100);

  sleep(2); // Slower pace for soak test
}
