# Helm Charts

This directory contains Helm charts for deploying CI/CD pipeline services to Kubernetes clusters.

## Purpose

Helm charts provide:
- Package management for Kubernetes applications
- Version control for deployments
- Environment-specific configuration
- Rollback capabilities

## Contents

### Service Charts

- **go-gateway/** - Helm chart for Go API Gateway
- **rust-store/** - Helm chart for Rust Data Store
- **python-telemetry/** - Helm chart for Python Telemetry Collector

## Usage

### Install Charts

```bash
# Install Go API Gateway
helm install go-api-gateway ./helm/go-gateway

# Install Rust Data Store
helm install rust-data-store ./helm/rust-store

# Install Python Telemetry
helm install python-telemetry ./helm/python-telemetry
```

### Upgrade Charts

```bash
# Upgrade Go API Gateway
helm upgrade go-api-gateway ./helm/go-gateway

# Upgrade Rust Data Store
helm upgrade rust-data-store ./helm/rust-store

# Upgrade Python Telemetry
helm upgrade python-telemetry ./helm/python-telemetry
```

### Uninstall Charts

```bash
# Uninstall Go API Gateway
helm uninstall go-api-gateway

# Uninstall Rust Data Store
helm uninstall rust-data-store

# Uninstall Python Telemetry
helm uninstall python-telemetry
```

## Chart Structure

Each chart follows standard Helm structure:

```
service-name/
├── Chart.yaml          # Chart metadata
├── values.yaml         # Default configuration values
├── templates/          # Kubernetes manifest templates
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── serviceaccount.yaml
│   └── _helpers.tpl    # Template helpers
└── README.md           # Chart documentation
```

## Configuration

### Values Override

Override default values with custom values file:

```bash
helm install go-api-gateway ./helm/go-gateway \
  --values custom-values.yaml
```

### Environment-Specific Values

Create environment-specific values files:

```yaml
# values-dev.yaml
replicaCount: 1
image:
  tag: dev

# values-prod.yaml
replicaCount: 3
image:
  tag: latest
resources:
  limits:
    cpu: 500m
    memory: 128Mi
```

### Common Configuration

**Replica Count:**
```yaml
replicaCount: 2
```

**Image Configuration:**
```yaml
image:
  repository: ghcr.io/username/ci-cd-pipeline/go-api-gateway
  pullPolicy: IfNotPresent
  tag: "latest"
```

**Resources:**
```yaml
resources:
  requests:
    cpu: 250m
    memory: 64Mi
  limits:
    cpu: 500m
    memory: 128Mi
```

**Service Configuration:**
```yaml
service:
  type: ClusterIP
  port: 8080
```

## Best Practices

### Chart Development

- Use semantic versioning
- Document chart changes
- Test charts locally
- Use Helm lint

### Configuration Management

- Use values files for configuration
- Separate environment-specific values
- Document configuration options
- Validate values schema

### Deployment

- Use Helm for production deployments
- Test in development first
- Use rollback for issues
- Monitor deployment status

## Troubleshooting

### Chart Installation Fails

Check:
- Kubernetes cluster connectivity
- Helm version compatibility
- Chart syntax errors
- Resource availability

### Values Not Applied

Verify:
- Values file syntax
- Template rendering
- Kubernetes resource limits
- Namespace configuration

### Upgrade Issues

Check:
- Chart version compatibility
- Breaking changes
- Resource conflicts
- Rollback availability

## Integration

### Argo CD

Helm charts integrate with Argo CD:
- GitOps deployment
- Automated sync
- Rollback support
- Configuration management

### CI/CD

Helm charts integrate with CI/CD:
- Automated deployments
- Environment-specific configuration
- Rollback capabilities
- Version control

## Documentation

- [Helm Documentation](https://helm.sh/docs/)
- [Helm Best Practices](https://helm.sh/docs/chart_best_practices/)
- [Chart Development Guide](https://helm.sh/docs/topics/charts/)
