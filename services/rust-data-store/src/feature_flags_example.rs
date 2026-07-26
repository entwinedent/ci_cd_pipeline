use feature_flags::{
    create_provider, Config, FeatureFlags, FLAG_ADVANCED_TELEMETRY, FLAG_EXPERIMENTAL_CACHE,
    FLAG_NEW_API_ENDPOINT, FLAG_RATE_LIMITING,
};

fn main() {
    // Initialize feature flag provider
    let config = Config::default();
    let provider = create_provider(config).expect(" Failed to create feature flag provider");

    // Example: Check feature flags in your application

    // Check if new API endpoint is enabled
    if provider.is_enabled(FLAG_NEW_API_ENDPOINT, false).unwrap_or(false) {
        println!("New API endpoint is enabled");
        // Use new endpoint implementation
    } else {
        println!("Using legacy API endpoint");
        // Use legacy implementation
    }

    // Check if advanced telemetry is enabled
    if provider.is_enabled(FLAG_ADVANCED_TELEMETRY, true).unwrap_or(true) {
        println!("Advanced telemetry is enabled");
        // Enable advanced telemetry collection
    }

    // Check if experimental cache is enabled
    if provider.is_enabled(FLAG_EXPERIMENTAL_CACHE, false).unwrap_or(false) {
        println!("Experimental cache is enabled");
        // Enable experimental caching layer
    }

    // Check if rate limiting is enabled
    if provider.is_enabled(FLAG_RATE_LIMITING, true).unwrap_or(true) {
        println!("Rate limiting is enabled");
        // Apply rate limiting
    }

    // Example: Get variant
    let variant = provider
        .get_variant("ui_theme", "default")
        .unwrap_or_else(|_| "default".to_string());
    println!("UI theme variant: {}", variant);

    // Environment variable override example
    // Set UNLEASH_FLAG_NEW_API_ENDPOINT=true to override
    // Set LD_FLAG_NEW_API_ENDPOINT=true for LaunchDarkly
}
