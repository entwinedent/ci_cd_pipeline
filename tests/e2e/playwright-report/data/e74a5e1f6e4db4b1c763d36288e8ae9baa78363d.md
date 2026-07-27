# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: api-tests.spec.ts >> Data Store API Operations >> DELETE /api/v1/data/{key} should delete data
- Location: api-tests.spec.ts:45:7

# Error details

```
Error: apiRequestContext.delete: socket hang up
Call log:
  - → DELETE http://localhost:8080/api/v1/data/test-key-1785121992521
    - user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.7922.34 Safari/537.36
    - accept: */*
    - accept-encoding: gzip,deflate,br

```

# Test source

```ts
  1   | import { test, expect } from '@playwright/test';
  2   | 
  3   | const API_GATEWAY_URL = process.env.API_GATEWAY_URL || 'http://localhost:8080';
  4   | const TELEMETRY_URL = process.env.TELEMETRY_URL || 'http://localhost:8000';
  5   | 
  6   | test.describe('API Gateway Health Checks', () => {
  7   |   test('GET /healthz should return healthy status', async ({ request }) => {
  8   |     const response = await request.get(`${API_GATEWAY_URL}/healthz`);
  9   |     expect(response.status()).toBe(200);
  10  |     const body = await response.json();
  11  |     expect(body.status).toBe('healthy');
  12  |   });
  13  | 
  14  |   test('GET /readyz should return ready status', async ({ request }) => {
  15  |     const response = await request.get(`${API_GATEWAY_URL}/readyz`);
  16  |     expect(response.status()).toBe(200);
  17  |     const body = await response.json();
  18  |     expect(body.status).toBe('ready');
  19  |   });
  20  | });
  21  | 
  22  | test.describe('Data Store API Operations', () => {
  23  |   const testKey = `test-key-${Date.now()}`;
  24  |   const testValue = Buffer.from('test-data-value');
  25  | 
  26  |   test('POST /api/v1/data should store data', async ({ request }) => {
  27  |     const response = await request.post(`${API_GATEWAY_URL}/api/v1/data`, {
  28  |       data: {
  29  |         key: testKey,
  30  |         value: testValue.toString('base64'),
  31  |       },
  32  |     });
  33  |     expect(response.status()).toBe(201);
  34  |     const body = await response.json();
  35  |     expect(body.status).toBe('success');
  36  |   });
  37  | 
  38  |   test('GET /api/v1/data/{key} should retrieve data', async ({ request }) => {
  39  |     const response = await request.get(`${API_GATEWAY_URL}/api/v1/data/${testKey}`);
  40  |     expect(response.status()).toBe(200);
  41  |     const body = await response.json();
  42  |     expect(body.key).toBe(testKey);
  43  |   });
  44  | 
  45  |   test('DELETE /api/v1/data/{key} should delete data', async ({ request }) => {
> 46  |     const response = await request.delete(`${API_GATEWAY_URL}/api/v1/data/${testKey}`);
      |                                          ^ Error: apiRequestContext.delete: socket hang up
  47  |     expect(response.status()).toBe(200);
  48  |     const body = await response.json();
  49  |     expect(body.status).toBe('success');
  50  |   });
  51  | });
  52  | 
  53  | test.describe('Telemetry Service Health Checks', () => {
  54  |   test('GET /healthz should return healthy status', async ({ request }) => {
  55  |     const response = await request.get(`${TELEMETRY_URL}/healthz`);
  56  |     expect(response.status()).toBe(200);
  57  |     const body = await response.json();
  58  |     expect(body.status).toBe('healthy');
  59  |   });
  60  | 
  61  |   test('GET /readyz should return ready status', async ({ request }) => {
  62  |     const response = await request.get(`${TELEMETRY_URL}/readyz`);
  63  |     expect(response.status()).toBe(200);
  64  |     const body = await response.json();
  65  |     expect(body.status).toBe('ready');
  66  |   });
  67  | });
  68  | 
  69  | test.describe('Telemetry Log Ingestion', () => {
  70  |   test('POST /api/v1/logs should ingest logs successfully', async ({ request }) => {
  71  |     const logEntry = {
  72  |       service: 'test-service',
  73  |       level: 'INFO',
  74  |       message: 'Test log message',
  75  |       timestamp: new Date().toISOString(),
  76  |       metadata: {
  77  |         test: true,
  78  |       },
  79  |     };
  80  | 
  81  |     const response = await request.post(`${TELEMETRY_URL}/api/v1/logs`, {
  82  |       data: logEntry,
  83  |     });
  84  |     expect(response.status()).toBe(200);
  85  |     const body = await response.json();
  86  |     expect(body.status).toBe('success');
  87  |   });
  88  | });
  89  | 
  90  | test.describe('Telemetry Metrics Ingestion', () => {
  91  |   test('POST /api/v1/metrics should ingest metrics successfully', async ({ request }) => {
  92  |     const metricData = {
  93  |       service: 'test-service',
  94  |       metric_name: 'test_metric',
  95  |       value: 42.5,
  96  |       timestamp: new Date().toISOString(),
  97  |       labels: {
  98  |         environment: 'test',
  99  |       },
  100 |     };
  101 | 
  102 |     const response = await request.post(`${TELEMETRY_URL}/api/v1/metrics`, {
  103 |       data: metricData,
  104 |     });
  105 |     expect(response.status()).toBe(200);
  106 |     const body = await response.json();
  107 |     expect(body.status).toBe('success');
  108 |   });
  109 | 
  110 |   test('GET /api/v1/anomalies should return recent anomalies', async ({ request }) => {
  111 |     const response = await request.get(`${TELEMETRY_URL}/api/v1/anomalies`);
  112 |     expect(response.status()).toBe(200);
  113 |     const body = await response.json();
  114 |     expect(body).toHaveProperty('anomalies');
  115 |     expect(Array.isArray(body.anomalies)).toBe(true);
  116 |   });
  117 | });
  118 | 
  119 | test.describe('Microservice Chain Integration', () => {
  120 |   const chainKey = `chain-test-${Date.now()}`;
  121 |   const chainValue = Buffer.from('chain-test-value');
  122 | 
  123 |   test('Complete flow: Gateway → Data Store → Telemetry', async ({ request }) => {
  124 |     // Step 1: Store data via Gateway
  125 |     const storeResponse = await request.post(`${API_GATEWAY_URL}/api/v1/data`, {
  126 |       data: {
  127 |         key: chainKey,
  128 |         value: chainValue.toString('base64'),
  129 |       },
  130 |     });
  131 |     expect(storeResponse.status()).toBe(201);
  132 | 
  133 |     // Step 2: Retrieve data via Gateway
  134 |     const getResponse = await request.get(`${API_GATEWAY_URL}/api/v1/data/${chainKey}`);
  135 |     expect(getResponse.status()).toBe(200);
  136 | 
  137 |     // Step 3: Send telemetry log about the operation
  138 |     const logResponse = await request.post(`${TELEMETRY_URL}/api/v1/logs`, {
  139 |       data: {
  140 |         service: 'api-gateway',
  141 |         level: 'INFO',
  142 |         message: `Data operation completed for key: ${chainKey}`,
  143 |         timestamp: new Date().toISOString(),
  144 |       },
  145 |     });
  146 |     expect(logResponse.status()).toBe(200);
```