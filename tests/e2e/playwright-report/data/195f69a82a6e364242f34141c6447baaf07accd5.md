# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: e2e-tests.spec.ts >> Service Availability Tests >> API Gateway is reachable from browser
- Location: e2e-tests.spec.ts:47:7

# Error details

```
Error: page.goto: net::ERR_EMPTY_RESPONSE at http://localhost:8080/healthz
Call log:
  - navigating to "http://localhost:8080/healthz", waiting until "load"

```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | const API_GATEWAY_URL = process.env.API_GATEWAY_URL || 'http://localhost:8080';
  4   | const TELEMETRY_URL = process.env.TELEMETRY_URL || 'http://localhost:8000';
  5   | 
  6   | test.describe('Browser E2E Tests - API Gateway', () => {
  7   |   test('Health check page loads and displays status', async ({ page }) => {
  8   |     await page.goto(`${API_GATEWAY_URL}/healthz`);
  9   |     const content = await page.content();
  10  |     expect(content).toContain('healthy');
  11  |   });
  12  | 
  13  |   test('Readiness check page loads and displays status', async ({ page }) => {
  14  |     await page.goto(`${API_GATEWAY_URL}/readyz`);
  15  |     const content = await page.content();
  16  |     expect(content).toContain('ready');
  17  |   });
  18  | 
  19  |   test('API responds within acceptable latency', async ({ page }) => {
  20  |     const startTime = Date.now();
  21  |     const response = await page.goto(`${API_GATEWAY_URL}/healthz`);
  22  |     const endTime = Date.now();
  23  |     const latency = endTime - startTime;
  24  |     
  25  |     expect(response?.status()).toBe(200);
  26  |     expect(latency).toBeLessThan(1000); // Less than 1 second
  27  |   });
  28  | });
  29  | 
  30  | test.describe('Browser E2E Tests - Telemetry Service', () => {
  31  |   test('Telemetry health check page loads', async ({ page }) => {
  32  |     await page.goto(`${TELEMETRY_URL}/healthz`);
  33  |     const content = await page.content();
  34  |     expect(content).toContain('healthy');
  35  |   });
  36  | 
  37  |   test('Telemetry anomalies endpoint returns JSON', async ({ page }) => {
  38  |     const response = await page.goto(`${TELEMETRY_URL}/api/v1/anomalies`);
  39  |     const content = await page.content();
  40  |     
  41  |     expect(response?.status()).toBe(200);
  42  |     expect(content).toContain('anomalies');
  43  |   });
  44  | });
  45  | 
  46  | test.describe('Service Availability Tests', () => {
  47  |   test('API Gateway is reachable from browser', async ({ page }) => {
> 48  |     const response = await page.goto(`${API_GATEWAY_URL}/healthz`);
      |                                 ^ Error: page.goto: net::ERR_EMPTY_RESPONSE at http://localhost:8080/healthz
  49  |     expect(response?.status()).toBe(200);
  50  |   });
  51  | 
  52  |   test('Telemetry service is reachable from browser', async ({ page }) => {
  53  |     const response = await page.goto(`${TELEMETRY_URL}/healthz`);
  54  |     expect(response?.status()).toBe(200);
  55  |   });
  56  | 
  57  |   test('Services respond concurrently', async ({ page }) => {
  58  |     const [gatewayResponse, telemetryResponse] = await Promise.all([
  59  |       page.goto(`${API_GATEWAY_URL}/healthz`),
  60  |       page.goto(`${TELEMETRY_URL}/healthz`),
  61  |     ]);
  62  | 
  63  |     expect(gatewayResponse?.status()).toBe(200);
  64  |     expect(telemetryResponse?.status()).toBe(200);
  65  |   });
  66  | });
  67  | 
  68  | test.describe('Error Handling Tests', () => {
  69  |   test('Non-existent endpoint returns 404', async ({ page }) => {
  70  |     const response = await page.goto(`${API_GATEWAY_URL}/non-existent`);
  71  |     expect(response?.status()).toBe(404);
  72  |   });
  73  | 
  74  |   test('Invalid data format returns error', async ({ page, request }) => {
  75  |     const response = await request.post(`${API_GATEWAY_URL}/api/v1/data`, {
  76  |       data: {
  77  |         invalid: 'data',
  78  |       },
  79  |     });
  80  |     expect(response.status()).toBeGreaterThanOrEqual(400);
  81  |   });
  82  | });
  83  | 
  84  | test.describe('Performance Tests', () => {
  85  |   test('Multiple concurrent requests complete successfully', async ({ page }) => {
  86  |     const requests = Array(10).fill(null).map((_, i) => 
  87  |       page.goto(`${API_GATEWAY_URL}/healthz`)
  88  |     );
  89  | 
  90  |     const responses = await Promise.all(requests);
  91  |     responses.forEach(response => {
  92  |       expect(response?.status()).toBe(200);
  93  |     });
  94  |   });
  95  | 
  96  |   test('Response time is consistent across multiple calls', async ({ page }) => {
  97  |     const latencies: number[] = [];
  98  |     
  99  |     for (let i = 0; i < 5; i++) {
  100 |       const start = Date.now();
  101 |       await page.goto(`${API_GATEWAY_URL}/healthz`);
  102 |       const end = Date.now();
  103 |       latencies.push(end - start);
  104 |     }
  105 | 
  106 |     const avgLatency = latencies.reduce((a, b) => a + b, 0) / latencies.length;
  107 |     expect(avgLatency).toBeLessThan(500); // Average less than 500ms
  108 |   });
  109 | });
  110 | 
```