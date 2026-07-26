package featureflags

import (
	"context"
	"fmt"
	"log"
)

// Example usage of feature flags in the Go API Gateway
func ExampleUsage() {
	// Initialize feature flag provider
	cfg := Config{
		Provider: ProviderUnleash, // or ProviderLaunchDarkly
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		log.Fatalf("Failed to create feature flag provider: %v", err)
	}
	defer provider.Close()

	// Create feature flags instance
	ff := NewFeatureFlags(provider)

	// Example: Check feature flags in your application
	ctx := context.Background()

	// Check if new API endpoint is enabled
	if ff.NewAPIEndpointEnabled(ctx) {
		fmt.Println("New API endpoint is enabled")
		// Use new endpoint implementation
	} else {
		fmt.Println("Using legacy API endpoint")
		// Use legacy implementation
	}

	// Check if advanced telemetry is enabled
	if ff.AdvancedTelemetryEnabled(ctx) {
		fmt.Println("Advanced telemetry is enabled")
		// Enable advanced telemetry collection
	}

	// Check if experimental cache is enabled
	if ff.ExperimentalCacheEnabled(ctx) {
		fmt.Println("Experimental cache is enabled")
		// Enable experimental caching layer
	}

	// Check if rate limiting is enabled
	if ff.RateLimitingEnabled(ctx) {
		fmt.Println("Rate limiting is enabled")
		// Apply rate limiting
	}

	// Example: Direct provider usage
	isEnabled, err := provider.IsEnabled(ctx, "custom_feature", false)
	if err != nil {
		log.Printf("Error checking custom feature: %v", err)
	} else if isEnabled {
		fmt.Println("Custom feature is enabled")
	}

	// Example: Get variant
	variant, err := provider.GetVariant(ctx, "ui_theme", "default")
	if err != nil {
		log.Printf("Error getting variant: %v", err)
	} else {
		fmt.Printf("UI theme variant: %s\n", variant)
	}
}
