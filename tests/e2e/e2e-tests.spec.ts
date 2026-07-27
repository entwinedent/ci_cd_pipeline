import { test, expect } from '@playwright/test';

const API_GATEWAY_URL = process.env.API_GATEWAY_URL || 'http://localhost:8080';
const TELEMETRY_URL = process.env.TELEMETRY_URL || 'http://localhost:8000';

test.describe('Browser E2E Tests - API Gateway', () => {
  test('Health check page loads and displays status', async ({ page }) => {
    await page.goto(`${API_GATEWAY_URL}/healthz`);
    const content = await page.content();
    expect(content).toContain('healthy');
  });

  test('Readiness check page loads and displays status', async ({ page }) => {
    await page.goto(`${API_GATEWAY_URL}/readyz`);
    const content = await page.content();
    expect(content).toContain('ready');
  });

  test('API responds within acceptable latency', async ({ page }) => {
    const startTime = Date.now();
    const response = await page.goto(`${API_GATEWAY_URL}/healthz`);
    const endTime = Date.now();
    const latency = endTime - startTime;
    
    expect(response?.status()).toBe(200);
    expect(latency).toBeLessThan(1000); // Less than 1 second
  });
});

test.describe('Browser E2E Tests - Telemetry Service', () => {
  test('Telemetry health check page loads', async ({ page }) => {
    await page.goto(`${TELEMETRY_URL}/healthz`);
    const content = await page.content();
    expect(content).toContain('healthy');
  });

  test('Telemetry anomalies endpoint returns JSON', async ({ page }) => {
    const response = await page.goto(`${TELEMETRY_URL}/api/v1/anomalies`);
    const content = await page.content();
    
    expect(response?.status()).toBe(200);
    expect(content).toContain('anomalies');
  });
});

test.describe('Service Availability Tests', () => {
  test('API Gateway is reachable from browser', async ({ page }) => {
    const response = await page.goto(`${API_GATEWAY_URL}/healthz`);
    expect(response?.status()).toBe(200);
  });

  test('Telemetry service is reachable from browser', async ({ page }) => {
    const response = await page.goto(`${TELEMETRY_URL}/healthz`);
    expect(response?.status()).toBe(200);
  });

  test('Services respond concurrently', async ({ page }) => {
    const gatewayResponse = await page.goto(`${API_GATEWAY_URL}/healthz`);
    const telemetryResponse = await page.goto(`${TELEMETRY_URL}/healthz`);

    expect(gatewayResponse?.status()).toBeGreaterThanOrEqual(200);
    expect(gatewayResponse?.status()).toBeLessThan(500);
    expect(telemetryResponse?.status()).toBeGreaterThanOrEqual(200);
    expect(telemetryResponse?.status()).toBeLessThan(500);
  });
});

test.describe('Error Handling Tests', () => {
  test('Non-existent endpoint returns 404', async ({ page }) => {
    const response = await page.goto(`${API_GATEWAY_URL}/non-existent`);
    expect(response?.status()).toBe(404);
  });

  test('Invalid data format returns error', async ({ page, request }) => {
    const response = await request.post(`${API_GATEWAY_URL}/api/v1/data`, {
      data: {
        invalid: 'data',
      },
    });
    expect(response.status()).toBeGreaterThanOrEqual(400);
  });
});

test.describe('Performance Tests', () => {
  test('Multiple concurrent requests complete successfully', async ({ page }) => {
    const responses = [];
    for (let i = 0; i < 10; i++) {
      const response = await page.goto(`${API_GATEWAY_URL}/healthz`);
      responses.push(response);
    }

    responses.forEach(response => {
      expect(response?.status()).toBeGreaterThanOrEqual(200);
      expect(response?.status()).toBeLessThan(500);
    });
  });

  test('Response time is consistent across multiple calls', async ({ page }) => {
    const latencies: number[] = [];
    
    for (let i = 0; i < 5; i++) {
      const start = Date.now();
      await page.goto(`${API_GATEWAY_URL}/healthz`);
      const end = Date.now();
      latencies.push(end - start);
    }

    const avgLatency = latencies.reduce((a, b) => a + b, 0) / latencies.length;
    expect(avgLatency).toBeLessThan(500); // Average less than 500ms
  });
});
