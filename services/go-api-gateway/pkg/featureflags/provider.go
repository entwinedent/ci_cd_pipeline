package featureflags

import (
	"context"
	"fmt"
	"os"
)

// ProviderType defines the feature flag provider
type ProviderType string

const (
	ProviderUnleash     ProviderType = "unleash"
	ProviderLaunchDarkly ProviderType = "launchdarkly"
)

// Provider is the interface for feature flag providers
type Provider interface {
	// IsEnabled checks if a feature flag is enabled
	IsEnabled(ctx context.Context, flagName string, defaultValue bool) (bool, error)
	// GetVariant returns the variant for a feature flag
	GetVariant(ctx context.Context, flagName string, defaultValue string) (string, error)
	// Close closes the provider connection
	Close() error
}

// Config holds the configuration for the feature flag provider
type Config struct {
	Provider   ProviderType
	Unleash    UnleashConfig
	LaunchDarkly LaunchDarklyConfig
}

// UnleashConfig holds Unleash-specific configuration
type UnleashConfig struct {
	URL       string
	APIToken  string
	AppName   string
	Environment string
}

// LaunchDarklyConfig holds LaunchDarkly-specific configuration
type LaunchDarklyConfig struct {
	SDKKey     string
	AppName    string
	Environment string
}

// NewProvider creates a new feature flag provider based on configuration
func NewProvider(cfg Config) (Provider, error) {
	// If provider not set, check environment variable
	if cfg.Provider == "" {
		providerEnv := os.Getenv("FEATURE_FLAG_PROVIDER")
		if providerEnv == "" {
			providerEnv = string(ProviderUnleash) // Default to Unleash
		}
		cfg.Provider = ProviderType(providerEnv)
	}

	switch cfg.Provider {
	case ProviderUnleash:
		return NewUnleashProvider(cfg.Unleash)
	case ProviderLaunchDarkly:
		return NewLaunchDarklyProvider(cfg.LaunchDarkly)
	default:
		return nil, fmt.Errorf("unknown feature flag provider: %s", cfg.Provider)
	}
}
