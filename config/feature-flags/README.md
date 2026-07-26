# Feature Flags Configuration

This directory contains centralized feature flag configuration for the CI/CD pipeline services.

## Purpose

Feature flags enable:
- Runtime feature toggling
- Gradual rollouts
- A/B testing
- Emergency feature disable
- Environment-specific behavior

## Contents

### Configuration File

- **config.yaml** - Central feature flag configuration defining:
  - Active feature flag provider (Unleash/LaunchDarkly)
  - Feature flag definitions
  - Default values
  - Environment overrides

## Usage

### Configuration Structure

```yaml
provider:
  type: "unleash"  # or "launchdarkly"
  unleash:
    url: "http://unleash.bleeding-edge.svc.cluster.local:4242"
    api_token: ""
  launchdarkly:
    sdk_key: ""

flags:
  new_api_endpoint:
    enabled: true
    default: false
    description: "Enable new API endpoint"
  
  experimental_caching:
    enabled: true
    default: false
    description: "Enable experimental caching"
```

### Service Integration

Services read configuration at startup:

```go
// Go
cfg := featureflags.Config{
    Provider: featureflags.ProviderUnleash,
}
provider, _ := featureflags.NewProvider(cfg)
```

```rust
// Rust
let config = FeatureFlagConfig {
    provider_type: ProviderType::Unleash,
};
let provider = create_provider(&config)?;
```

```python
# Python
config = FeatureFlagConfig(
    provider_type=ProviderType.UNLEASH,
)
provider = create_provider(config)
```

### Environment Overrides

Override configuration with environment variables:

```bash
export FEATURE_FLAG_PROVIDER=launchdarkly
export LAUNCHDARKLY_SDK_KEY=your-sdk-key
export UNLEASH_URL=http://unleash.local:4242
```

## Feature Flag Providers

### Unleash (Self-Hosted)

**Advantages:**
- Self-hosted, no external dependencies
- Full control over data
- No cost for self-hosting
- Customizable

**Configuration:**
- Deploy Unleash server
- Configure feature flags in UI
- Use SDK for integration

### LaunchDarkly (Commercial)

**Advantages:**
- Managed service
- Advanced features
- Analytics and insights
- Enterprise support

**Configuration:**
- Create LaunchDarkly account
- Configure feature flags in dashboard
- Use SDK key for integration

## Best Practices

### Flag Design

- Use descriptive flag names
- Document flag purpose
- Set sensible defaults
- Plan flag lifecycle

### Flag Management

- Regularly review active flags
- Remove unused flags
- Document flag changes
- Monitor flag usage

### Testing

- Test with flags enabled/disabled
- Test flag transitions
- Test default values
- Test provider switching

## Troubleshooting

### Provider Connection Issues

Check:
- Provider URL configuration
- Network connectivity
- API credentials
- Service availability

### Flag Not Working

Verify:
- Flag name matches configuration
- Provider is accessible
- SDK initialization
- Environment variables

### Switching Providers

Steps:
1. Update configuration file
2. Set environment variables
3. Restart services
4. Verify flag functionality

## Documentation

- [Unleash Documentation](https://docs.unleash.io/)
- [LaunchDarkly Documentation](https://docs.launchdarkly.com/)
- [Feature Flag Best Practices](https://martinfowler.com/articles/feature-toggles.html)
