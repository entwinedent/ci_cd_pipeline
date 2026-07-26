# Feature Flag Abstraction Layer

This directory contains the implementation of the feature flag abstraction layer for the CI/CD pipeline, supporting both Unleash (self-hosted) and LaunchDarkly (commercial) providers.

## Architecture

The feature flag system consists of:

1. **Abstraction Layer**: Language-specific interfaces that provide a unified API for feature flag operations
2. **Provider Implementations**: Concrete implementations for Unleash and LaunchDarkly
3. **Provider Switching**: Runtime configuration to switch between providers without code changes
4. **Kubernetes Deployment**: Self-hosted Unleash deployment with PostgreSQL

## Components

### Unleash Self-Hosted

- `namespace.yaml` - Kubernetes namespace for Unleash
- `postgres-secret.yaml` - PostgreSQL credentials
- `postgres-deployment.yaml` - PostgreSQL database deployment
- `unleash-config.yaml` - Unleash server configuration
- `unleash-deployment.yaml` - Unleash server deployment with ingress
- `launchdarkly-secret.yaml` - LaunchDarkly SDK key secret
- `feature-flag-configmap.yaml` - Feature flag provider configuration

### Service Implementations

#### Go (`services/go-api-gateway/pkg/featureflags/`)
- `provider.go` - Provider interface and factory
- `unleash.go` - Unleash provider implementation
- `launchdarkly.go` - LaunchDarkly provider implementation
- `flags.go` - Convenient feature flag wrapper

#### Rust (`services/rust-data-store/src/feature_flags/`)
- `mod.rs` - Module definition with provider trait
- `unleash.rs` - Unleash provider implementation
- `launchdarkly.rs` - LaunchDarkly provider implementation

#### Python (`services/python-telemetry/src/feature_flags/`)
- `__init__.py` - Module definition with provider protocol
- `unleash.py` - Unleash provider implementation
- `launchdarkly.py` - LaunchDarkly provider implementation

## Deployment

### Deploy Unleash Self-Hosted

```bash
# Deploy to Kubernetes
kubectl apply -f manifests/bleeding-edge/unleash/namespace.yaml
kubectl apply -f manifests/bleeding-edge/unleash/postgres-secret.yaml
kubectl apply -f manifests/bleeding-edge/unleash/postgres-deployment.yaml
kubectl apply -f manifests/bleeding-edge/unleash/unleash-config.yaml
kubectl apply -f manifests/bleeding-edge/unleash/unleash-deployment.yaml
kubectl apply -f manifests/bleeding-edge/unleash/feature-flag-configmap.yaml
```

### Access Unleash UI

```bash
# Port forward to access Unleash UI
kubectl port-forward -n unleash svc/unleash 4242:4242

# Open browser to http://localhost:4242
# Default credentials: unleash / unleash (change in production)
```

## Configuration

### Provider Selection

Set the `FEATURE_FLAG_PROVIDER` environment variable to switch providers:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: feature-flag-config
  namespace: default
data:
  FEATURE_FLAG_PROVIDER: "unleash"  # or "launchdarkly"
```

### Unleash Configuration

```yaml
UNLEASH_URL: "http://unleash.unleash.svc.cluster.local:4242"
UNLEASH_API_TOKEN: "default:development.unleash-insecure-api-token-change-me"
UNLEASH_APP_NAME: "ci-cd-pipeline"
UNLEASH_ENVIRONMENT: "development"
```

### LaunchDarkly Configuration

```yaml
LAUNCHDARKLY_SDK_KEY: "your-sdk-key-here"
LAUNCHDARKLY_APP_NAME: "ci-cd-pipeline"
LAUNCHDARKLY_ENVIRONMENT: "development"
```

## Usage Examples

### Go

```go
cfg := featureflags.Config{
    Provider: featureflags.ProviderUnleash,
}
provider, _ := featureflags.NewProvider(cfg)
ff := featureflags.NewFeatureFlags(provider)

if ff.NewAPIEndpointEnabled(ctx) {
    // Use new implementation
}
```

### Rust

```rust
let config = Config::default();
let provider = create_provider(config).unwrap();

if provider.is_enabled(FLAG_NEW_API_ENDPOINT, false).unwrap_or(false) {
    // Use new implementation
}
```

### Python

```python
config = Config()
provider = create_provider(config)

if provider.is_enabled(FLAG_NEW_API_ENDPOINT, False):
    # Use new implementation
```

## Feature Flags

Predefined feature flags available across all services:

- `new_api_endpoint` - Enable new API endpoint implementation (default: false)
- `advanced_telemetry` - Enable advanced telemetry collection (default: true)
- `experimental_cache` - Enable experimental caching layer (default: false)
- `rate_limiting` - Enable rate limiting on API gateway (default: true)

## Environment Variable Overrides

For local development, you can override feature flags using environment variables:

### Unleash
```bash
export UNLEASH_FLAG_NEW_API_ENDPOINT=true
export UNLEASH_VARIANT_UI_THEME=dark
```

### LaunchDarkly
```bash
export LD_FLAG_NEW_API_ENDPOINT=true
export LD_VARIANT_UI_THEME=dark
```

## Production Considerations

1. **Security**: Change default passwords and API tokens in production
2. **Persistence**: Use persistent volumes for PostgreSQL in production
3. **High Availability**: Deploy Unleash with multiple replicas in production
4. **Backup**: Implement regular backups of the Unleash database
5. **Monitoring**: Add monitoring and alerting for Unleash server health
6. **Rate Limiting**: Configure rate limiting on the Unleash API

## Migration from Unleash to LaunchDarkly

To switch from Unleash to LaunchDarkly:

1. Set `FEATURE_FLAG_PROVIDER=launchdarkly` in the ConfigMap
2. Add `LAUNCHDARKLY_SDK_KEY` to the LaunchDarkly secret
3. Redeploy services
4. No code changes required due to abstraction layer

## Uninstallation

```bash
kubectl delete -f manifests/bleeding-edge/unleash/
kubectl delete namespace unleash
```
