# OpenCost Deployment for FinOps

This directory contains the OpenCost deployment configuration for runtime cost observability and FinOps governance.

## Architecture

OpenCost provides:
- **Real-time Cost Monitoring**: CPU, memory, storage, and network costs
- **Prometheus Integration**: Export cost metrics to Prometheus
- **Multi-Cloud Support**: AWS, GCP, Azure pricing models
- **Namespace/Workload Allocation**: Cost breakdown by namespace, deployment, and pod

## Components

- **opencost-deployment.yaml**: OpenCost deployment with RBAC
- **opencost-secret.yaml**: Cloud provider API key secret
- **opencost-pricing ConfigMap**: Custom pricing configuration
- **grafana-dashboards**: FinOps dashboards for Grafana

## Deployment

```bash
# Deploy OpenCost
kubectl apply -f manifests/bleeding-edge/finops/opencost/

# Verify deployment
kubectl get pods -n opencost
kubectl get svc -n opencost

# Port forward to access OpenCost UI
kubectl port-forward -n opencost svc/opencost 9003:9003
```

## Configuration

### Custom Pricing

Edit the `opencost-pricing` ConfigMap to customize pricing:

```yaml
data:
  pricing.json: |
    {
      "CPU": "0.0316",           # $0.0316 per vCPU hour
      "RAM": "0.00423",          # $0.00423 per GB hour
      "GPU": "0.95",             # $0.95 per GPU hour
      "Storage": "0.0001",       # $0.0001 per GB hour
      "InternetNetworkEgress": "0.12"  # $0.12 per GB
    }
```

### Prometheus Integration

OpenCost scrapes Prometheus metrics and exports cost data:

```yaml
env:
- name: PROMETHEUS_ENDPOINT
  value: "http://prometheus-server.monitoring.svc.cluster.local:9090"
```

## Grafana Dashboards

Import the FinOps dashboard:

1. Open Grafana
2. Navigate to Dashboards → Import
3. Upload `finops-dashboard.json`
4. Select Prometheus data source

## Cost Optimization

### Resource Rightsizing

Monitor resource efficiency and rightsize based on actual usage:

```bash
# Check CPU efficiency
kubectl top pods

# Compare with allocated resources
kubectl get pods -o jsonpath='{.spec.containers[*].resources}'
```

### Idle Resource Detection

Identify and remove idle resources:

```bash
# Find deployments with low CPU usage
kubectl get deployments --all-namespaces -o json | jq '.items[] | select(.spec.replicas > 0)'
```

### Budget Alerts

Set up Prometheus alerts for cost thresholds:

```yaml
groups:
- name: finops
  rules:
  - alert: HighMonthlyCost
    expr: sum(container_cpu_allocation) * 30 * 24 > 1000
    annotations:
      summary: "Monthly cost exceeds $1000"
```

## KEDA Integration

Drive autoscaling based on cost metrics:

```yaml
triggers:
- type: prometheus
  metadata:
    serverAddress: http://prometheus-server.monitoring.svc.cluster.local:9090
    metricName: cost_per_request
    threshold: "0.01"
    query: "sum(rate(container_cpu_usage_seconds_total[5m])) / sum(rate(http_requests_total[5m]))"
```

## Testing

### Mock Pricing

For testing without actual cloud costs, use the mock pricing in the ConfigMap.

### Cost Prediction

Use the mock kubectl-cost tool:

```bash
python tools/kubectl-cost-mock/kubectl_cost_mock.py --path k8s/base
```

## Security

- **RBAC**: OpenCost has minimal RBAC permissions
- **Service Account**: Dedicated service account for cost monitoring
- **Network Policies**: Restrict access to OpenCost service
- **Secrets**: Cloud provider API keys stored in Kubernetes secrets

## Monitoring

Monitor OpenCost health:

```bash
# Check OpenCost pod status
kubectl get pods -n opencost -l app=opencost

# Check OpenCost logs
kubectl logs -n opencost -l app=opencost

# Check Prometheus metrics
curl http://opencost.opencost.svc.cluster.local:9004/metrics
```

## Troubleshooting

### OpenCost Not Reporting Costs

1. Check Prometheus endpoint configuration
2. Verify RBAC permissions
3. Check OpenCost logs for errors
4. Ensure pricing ConfigMap is valid

### High Cost Alerts

1. Review resource allocations in Kubernetes manifests
2. Check for runaway pods or deployments
3. Verify pricing configuration is accurate
4. Consider using spot instances for non-critical workloads
