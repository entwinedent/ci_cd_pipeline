import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');
const API_GATEWAY_URL = __ENV.API_GATEWAY_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '30s', target: 100 },  // Ramp to 100
    { duration: '30s', target: 200 },  // Ramp to 200
    { duration: '30s', target: 500 },  // Ramp to 500
    { duration: '1m', target: 500 },   // Sustained extreme load
    { duration: '30s', target: 1000 }, // Push to 1000
    { duration: '30s', target: 1000 }, // Sustained extreme
    { duration: '30s', target: 0 },    // Crash recovery
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // Very lenient during stress
    errors: ['rate<0.2'], // Allow higher errors during stress
  },
};

export default function () {
  // Minimal checks during stress test - focus on system stability
  const healthRes = http.get(`${API_GATEWAY_URL}/healthz`, {
    timeout: '5s',
  });
  
  check(healthRes, {
    'health check status 200': (r) => r.status === 200,
    'health check timeout': (r) => r.timings.duration < 5000,
  }) || errorRate.add(1);

  // Rapid data operations
  const testKey = `stress-${__VU}-${Date.now()}`;
  http.post(`${API_GATEWAY_URL}/api/v1/data`, JSON.stringify({
    key: testKey,
    value: 'x'.repeat(1000), // Larger payload
  }), {
    headers: { 'Content-Type': 'application/json' },
    timeout: '5s',
  });

  sleep(0.05); // Very rapid requests
}
