# Python Telemetry Collector

AI-powered anomaly detection and log analysis service with OpenTelemetry integration and FastAPI endpoints.

## Architecture

The Telemetry Collector provides:
- Real-time log ingestion and processing
- AI-powered anomaly detection
- OpenTelemetry span processing
- FastAPI REST endpoints
- Auto-remediation webhook triggers

## OpenTelemetry Integration

### Span Processing

The collector processes OpenTelemetry spans for distributed tracing:

**Features**
- Span aggregation and correlation
- Latency analysis and percentile calculation
- Error rate monitoring
- Service dependency mapping

**Configuration**
```python
from opentelemetry import trace
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter

exporter = OTLPSpanExporter(endpoint="http://tempo:4317")
trace.get_tracer_provider().add_span_processor(
    BatchSpanProcessor(exporter)
)
```

### Metrics Collection

Prometheus metrics for observability:
- Request latency histogram
- Error rate counter
- Active requests gauge
- Anomaly detection metrics

## AI Anomaly Detection

### Algorithm

Uses statistical analysis for anomaly detection:

**Methods**
1. **Z-Score Analysis**: Detect outliers based on standard deviation
2. **Moving Average**: Identify trends and patterns
3. **Percentile Analysis**: Detect latency spikes
4. **Rate of Change**: Identify sudden metric changes

**Configuration**
- Sensitivity threshold: 3 standard deviations
- Window size: 5 minutes
- Minimum samples: 100

### Auto-Remediation

Triggers rollback on detected anomalies:

**Triggers**
- Error rate > 5% for 5 minutes
- P99 latency > 1 second for 3 minutes
- Anomaly score > threshold

**Actions**
- Argo CD webhook for rollback
- Slack notification
- Incident creation

## FastAPI Endpoints

### API Endpoints

- `GET /healthz` - Liveness probe
- `GET /readyz` - Readiness probe
- `POST /api/v1/logs` - Ingest logs
- `GET /api/v1/metrics` - Query metrics
- `GET /api/v1/anomalies` - Get detected anomalies

### Log Ingestion

```python
POST /api/v1/logs
Content-Type: application/json

{
  "service": "api-gateway",
  "level": "info",
  "message": "Request processed",
  "timestamp": "2024-01-01T00:00:00Z",
  "metadata": {
    "request_id": "abc-123",
    "latency_ms": 50
  }
}
```

### Metrics Query

```python
GET /api/v1/metrics?service=api-gateway&time_range=1h

{
  "service": "api-gateway",
  "time_range": "1h",
  "metrics": {
    "request_count": 10000,
    "error_rate": 0.01,
    "p50_latency_ms": 45,
    "p95_latency_ms": 85,
    "p99_latency_ms": 120
  }
}
```

## Configuration

### Environment Variables
- `PORT` - Service listening port (default: 8000)
- `LOG_LEVEL` - Logging level (default: `info`)
- `ARGOCD_WEBHOOK_URL` - Argo CD webhook for auto-remediation
- `SLACK_WEBHOOK_URL` - Slack webhook for notifications
- `ANOMALY_THRESHOLD` - Z-score threshold (default: 3.0)
- `WINDOW_SIZE_MINUTES` - Analysis window (default: 5)

### Service Specifications
- **Port**: 8000
- **Health**: `/healthz` (liveness), `/readyz` (readiness)
- **Metrics**: Prometheus metrics on `/metrics`

## Development

### Setup
```bash
cd services/python-telemetry
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```

### Run
```bash
uvicorn src.main:app --reload --port 8000
```

### Test
```bash
pytest src/
pytest --cov=src --cov-report=html
```

### Docker Build
```bash
docker build -t python-telemetry:latest .
```

## Performance

- **Image Size**: < 100MB
- **Response Time**: < 50ms (p95)
- **Throughput**: > 1,000 logs/sec
- **Memory**: < 200MB

## Dependencies

- `fastapi` - Web framework
- `uvicorn` - ASGI server
- `pydantic` - Data validation
- `httpx` - HTTP client
- `opentelemetry-api` - OpenTelemetry API
- `opentelemetry-sdk` - OpenTelemetry SDK
- `numpy` - Statistical analysis
- `scikit-learn` - Machine learning

## Security

- Input validation with Pydantic
- Rate limiting on all endpoints
- Sensitive data redaction in logs
- SPIFFE/SPIRE mTLS for service communication
