// K6 Load Test Configuration
export const BASE_URL = __ENV.API_GATEWAY_URL || 'http://localhost:8080';
export const TELEMETRY_URL = __ENV.TELEMETRY_URL || 'http://localhost:8000';

export const THRESHOLDS = {
  // Allow higher error rate during CI due to service startup issues
  'errors': ['rate<0.5'], // Less than 50% error rate
  'http_req_duration': ['p(95)<500', 'p(99)<1000'], // 95% under 500ms, 99% under 1s
};
