package featureflags

import "context"

// Feature flag names
const (
	FlagNewAPIEndpoint     = "new_api_endpoint"
	FlagAdvancedTelemetry  = "advanced_telemetry"
	FlagExperimentalCache  = "experimental_cache"
	FlagRateLimiting       = "rate_limiting"
)

// FeatureFlags provides a convenient interface to check feature flags
type FeatureFlags struct {
	provider Provider
}

// NewFeatureFlags creates a new FeatureFlags instance
func NewFeatureFlags(provider Provider) *FeatureFlags {
	return &FeatureFlags{
		provider: provider,
	}
}

// NewAPIEndpointEnabled checks if the new API endpoint is enabled
func (ff *FeatureFlags) NewAPIEndpointEnabled(ctx context.Context) bool {
	enabled, err := ff.provider.IsEnabled(ctx, FlagNewAPIEndpoint, false)
	if err != nil {
		return false
	}
	return enabled
}

// AdvancedTelemetryEnabled checks if advanced telemetry is enabled
func (ff *FeatureFlags) AdvancedTelemetryEnabled(ctx context.Context) bool {
	enabled, err := ff.provider.IsEnabled(ctx, FlagAdvancedTelemetry, true)
	if err != nil {
		return true
	}
	return enabled
}

// ExperimentalCacheEnabled checks if experimental cache is enabled
func (ff *FeatureFlags) ExperimentalCacheEnabled(ctx context.Context) bool {
	enabled, err := ff.provider.IsEnabled(ctx, FlagExperimentalCache, false)
	if err != nil {
		return false
	}
	return enabled
}

// RateLimitingEnabled checks if rate limiting is enabled
func (ff *FeatureFlags) RateLimitingEnabled(ctx context.Context) bool {
	enabled, err := ff.provider.IsEnabled(ctx, FlagRateLimiting, true)
	if err != nil {
		return true
	}
	return enabled
}
