"""Example usage of feature flags in Python telemetry service."""

from src.feature_flags import (
    Config,
    create_provider,
    FLAG_ADVANCED_TELEMETRY,
    FLAG_EXPERIMENTAL_CACHE,
    FLAG_NEW_API_ENDPOINT,
    FLAG_RATE_LIMITING,
)


def main():
    # Initialize feature flag provider
    config = Config()
    provider = create_provider(config)

    try:
        # Example: Check feature flags in your application

        # Check if new API endpoint is enabled
        if provider.is_enabled(FLAG_NEW_API_ENDPOINT, False):
            print("New API endpoint is enabled")
            # Use new endpoint implementation
        else:
            print("Using legacy API endpoint")
            # Use legacy implementation

        # Check if advanced telemetry is enabled
        if provider.is_enabled(FLAG_ADVANCED_TELEMETRY, True):
            print("Advanced telemetry is enabled")
            # Enable advanced telemetry collection

        # Check if experimental cache is enabled
        if provider.is_enabled(FLAG_EXPERIMENTAL_CACHE, False):
            print("Experimental cache is enabled")
            # Enable experimental caching layer

        # Check if rate limiting is enabled
        if provider.is_enabled(FLAG_RATE_LIMITING, True):
            print("Rate limiting is enabled")
            # Apply rate limiting

        # Example: Get variant
        variant = provider.get_variant("ui_theme", "default")
        print(f"UI theme variant: {variant}")

        # Environment variable override example
        # Set UNLEASH_FLAG_NEW_API_ENDPOINT=true to override
        # Set LD_FLAG_NEW_API_ENDPOINT=true for LaunchDarkly

    finally:
        provider.close()


if __name__ == "__main__":
    main()
