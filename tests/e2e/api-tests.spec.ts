import { test, expect } from '@playwright/test';

const API_GATEWAY_URL = process.env.API_GATEWAY_URL || 'http://localhost:8080';
const TELEMETRY_URL = process.env.TELEMETRY_URL || 'http://localhost:8000';

test.describe('API Gateway Health Checks', () => {
  test('GET /healthz should return healthy status', async ({ request }) => {
    const response = await request.get(`${API_GATEWAY_URL}/healthz`);
    expect(response.status()).toBe(200);
  });

  test('GET /readyz should return ready status', async ({ request }) => {
    const response = await request.get(`${API_GATEWAY_URL}/readyz`);
    expect(response.status()).toBe(200);
  });
});

test.describe('Data Store API Operations', () => {
  const testKey = `test-key-${Date.now()}`;
  const testValue = Buffer.from('test-data-value');

  test('POST /api/v1/data should store data', async ({ request }) => {
    const response = await request.post(`${API_GATEWAY_URL}/api/v1/data/${testKey}`, {
      data: {
        value: testValue.toString('base64'),
        ttl_seconds: 3600,
      },
    });
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('success');
  });

  test('GET /api/v1/data/{key} should retrieve data', async ({ request }) => {
    const response = await request.get(`${API_GATEWAY_URL}/api/v1/data/${testKey}`);
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.key).toBe(testKey);
  });

  test('DELETE /api/v1/data/{key} should delete data', async ({ request }) => {
    const response = await request.delete(`${API_GATEWAY_URL}/api/v1/data/${testKey}`);
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('success');
  });
});

test.describe('Telemetry Service Health Checks', () => {
  test('GET /healthz should return healthy status', async ({ request }) => {
    const response = await request.get(`${TELEMETRY_URL}/healthz`);
    expect(response.status()).toBe(200);
  });

  test('GET /readyz should return ready status', async ({ request }) => {
    const response = await request.get(`${TELEMETRY_URL}/readyz`);
    expect(response.status()).toBe(200);
  });
});

test.describe('Telemetry Log Ingestion', () => {
  test('POST /api/v1/logs should ingest logs successfully', async ({ request }) => {
    const logEntry = {
      service: 'test-service',
      level: 'INFO',
      message: 'Test log message',
      timestamp: new Date().toISOString(),
      metadata: {
        test: true,
      },
    };

    const response = await request.post(`${TELEMETRY_URL}/api/v1/logs`, {
      data: logEntry,
    });
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('success');
  });
});

test.describe('Telemetry Metrics Ingestion', () => {
  test('POST /api/v1/metrics should ingest metrics successfully', async ({ request }) => {
    const metricData = {
      service: 'test-service',
      metric_name: 'test_metric',
      value: 42.5,
      timestamp: new Date().toISOString(),
      labels: {
        environment: 'test',
      },
    };

    const response = await request.post(`${TELEMETRY_URL}/api/v1/metrics`, {
      data: metricData,
    });
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body.status).toBe('success');
  });

  test('GET /api/v1/anomalies should return recent anomalies', async ({ request }) => {
    const response = await request.get(`${TELEMETRY_URL}/api/v1/anomalies`);
    expect(response.status()).toBe(200);
    const body = await response.json();
    expect(body).toHaveProperty('anomalies');
    expect(Array.isArray(body.anomalies)).toBe(true);
  });
});

test.describe('Microservice Chain Integration', () => {
  const chainKey = `chain-test-${Date.now()}`;
  const chainValue = Buffer.from('chain-test-value');

  test('Complete flow: Gateway → Data Store → Telemetry', async ({ request }) => {
    // Step 1: Store data via Gateway
    const storeResponse = await request.post(`${API_GATEWAY_URL}/api/v1/data/${chainKey}`, {
      data: {
        value: chainValue.toString('base64'),
        ttl_seconds: 3600,
      },
    });
    expect(storeResponse.status()).toBe(200);

    // Step 2: Retrieve data via Gateway
    const getResponse = await request.get(`${API_GATEWAY_URL}/api/v1/data/${chainKey}`);
    expect(getResponse.status()).toBe(200);

    // Step 3: Send telemetry log about the operation
    const logResponse = await request.post(`${TELEMETRY_URL}/api/v1/logs`, {
      data: {
        service: 'api-gateway',
        level: 'INFO',
        message: `Data operation completed for key: ${chainKey}`,
        timestamp: new Date().toISOString(),
      },
    });
    expect(logResponse.status()).toBe(200);

    // Step 4: Send metric about the operation
    const metricResponse = await request.post(`${TELEMETRY_URL}/api/v1/metrics`, {
      data: {
        service: 'api-gateway',
        metric_name: 'data_operations',
        value: 1,
        timestamp: new Date().toISOString(),
      },
    });
    expect(metricResponse.status()).toBe(200);

    // Cleanup
    await request.delete(`${API_GATEWAY_URL}/api/v1/data/${chainKey}`);
  });
});
