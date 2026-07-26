# Prometheus Configuration

This directory contains Prometheus monitoring configuration for the CI/CD pipeline services.

## Purpose

Prometheus provides:
- Metrics collection and storage
- Alerting and notification
- Service monitoring
- Performance analysis

## Contents

### Installation

- **install-prometheus.sh** - Script to install Prometheus in Kubernetes cluster
- **values.yaml** - Helm values for Prometheus configuration

## Usage

### Installation

```bash
# Install Prometheus
bash config/prometheus/install-prometheus.sh

# Verify installation
kubectl get pods -n monitoring
kubectl get svc -n monitoring
```

### Access Prometheus UI

```bash
# Port forward Prometheus
kubectl port-forward svc/prometheus-server -n monitoring 9090:9090

# Access at http://localhost:9090
```

## Configuration

### Scrape Configurations

Prometheus scrapes metrics from:
- Go API Gateway: `http://go-api-gateway:8080/metrics`
- Rust Data Store: `http://rust-data-store:50051/metrics`
- Python Telemetry: `http://python-telemetry:8000/metrics`
- Kubernetes nodes and pods
- Cilium networking metrics

### Alerting Rules

Configure alerting for:
- High error rates
- High latency
- Resource exhaustion
- Service unavailability

### Recording Rules

Pre-compute frequently used queries:
- Request rate by service
- Error rate by service
- Latency percentiles
- Resource utilization

## Metrics

### Application Metrics

**Go API Gateway:**
- `http_requests_total` - Total HTTP requests
- `http_request_duration_seconds` - Request latency
- `grpc_client_requests_total` - gRPC requests
- `feature_flag_checks_total` - Feature flag checks

**Rust Data Store:**
- `grpc_server_requests_total` - gRPC requests
- `grpc_server_latency_seconds` - Request latency
- `store_operations_total` - Store operations
- `cache_hits_total` - Cache hits

**Python Telemetry:**
- `log_ingestion_total` - Logs ingested
- `anomaly_detection_total` - Anomalies detected
- `metric_queries_total` - Metric queries
- `webhook_calls_total` - Webhook calls

### System Metrics

- CPU utilization
- Memory usage
- Network traffic
- Disk I/O
- Container metrics

## Best Practices

### Metric Design

- Use consistent naming conventions
- Include relevant labels
- Document metric purposes
- Avoid high cardinality labels

### Alerting

- Set appropriate thresholds
- Avoid alert fatigue
- Use severity levels
- Include runbooks

### Performance

- Optimize scrape intervals
- Use recording rules
- Configure retention policies
- Monitor Prometheus performance

## Troubleshooting

### Metrics Not Appearing

Check:
- Service metrics endpoint
- Prometheus scrape configuration
- Network connectivity
- Service discovery

### High Memory Usage

Verify:
- Scrape interval configuration
- Number of high cardinality metrics
- Retention policy
- Recording rules

### Alerting Not Working

Check:
- Alert rule configuration
- Alertmanager configuration
- Notification channels
- Alert evaluation interval

## Integration

### Grafana

Prometheus integrates with Grafana for visualization:
- Pre-built dashboards
- Custom dashboards
- Alert visualization
- Historical analysis

### Alertmanager

Prometheus uses Alertmanager for:
- Alert routing
- Notification management
- Alert grouping
- Silencing and inhibition

## Documentation

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Prometheus Best Practices](https://prometheus.io/docs/practices/)
- [Grafana Integration](https://grafana.com/docs/grafana/latest/datasources/prometheus/)
