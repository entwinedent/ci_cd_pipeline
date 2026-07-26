import { Pact } from '@pact-foundation/pact';
import { like, regex } from '@pact-foundation/pact/dsl/matchers';
import path from 'path';

const provider = new Pact({
  consumer: 'go-api-gateway',
  provider: 'rust-data-store',
  port: 50051,
  log: path.resolve(process.cwd(), 'logs', 'pact.log'),
  dir: path.resolve(process.cwd(), 'pacts'),
  logLevel: 'INFO',
});

describe('Go API Gateway Consumer Contract Tests', () => {
  beforeAll(async () => {
    await provider.setup();
  });

  afterAll(async () => {
    await provider.finalize();
  });

  afterEach(async () => {
    await provider.verify();
  });

  describe('Data Store Operations', () => {
    test('Set data operation', async () => {
      await provider.addInteraction({
        state: 'data store is available',
        uponReceiving: 'a request to set data',
        withRequest: {
          method: 'POST',
          path: '/v1/data',
          headers: {
            'Content-Type': 'application/grpc',
          },
          body: {
            key: like('test-key'),
            value: like(Buffer.from('test-value')),
            ttl_seconds: like(3600),
          },
        },
        willRespondWith: {
          status: 200,
          headers: {
            'Content-Type': 'application/grpc',
          },
          body: {
            success: true,
            message: like('Value stored successfully'),
          },
        },
      });
    });

    test('Get data operation', async () => {
      await provider.addInteraction({
        state: 'data exists for key',
        uponReceiving: 'a request to get data',
        withRequest: {
          method: 'GET',
          path: regex('/v1/data/[a-zA-Z0-9-]+', '/v1/data/test-key'),
          headers: {
            'Content-Type': 'application/grpc',
          },
        },
        willRespondWith: {
          status: 200,
          headers: {
            'Content-Type': 'application/grpc',
          },
          body: {
            value: like(Buffer.from('test-value')),
            found: true,
            message: like('Value retrieved successfully'),
          },
        },
      });
    });

    test('Delete data operation', async () => {
      await provider.addInteraction({
        state: 'data exists for key',
        uponReceiving: 'a request to delete data',
        withRequest: {
          method: 'DELETE',
          path: regex('/v1/data/[a-zA-Z0-9-]+', '/v1/data/test-key'),
          headers: {
            'Content-Type': 'application/grpc',
          },
        },
        willRespondWith: {
          status: 200,
          headers: {
            'Content-Type': 'application/grpc',
          },
          body: {
            success: true,
            message: like('Value deleted successfully'),
          },
        },
      });
    });

    test('Health check operation', async () => {
      await provider.addInteraction({
        state: 'service is healthy',
        uponReceiving: 'a health check request',
        withRequest: {
          method: 'GET',
          path: '/v1/health',
          headers: {
            'Content-Type': 'application/grpc',
          },
        },
        willRespondWith: {
          status: 200,
          headers: {
            'Content-Type': 'application/grpc',
          },
          body: {
            healthy: true,
            message: like('Service is healthy'),
          },
        },
      });
    });
  });
});
