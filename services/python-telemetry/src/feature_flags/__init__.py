"""Feature flag abstraction layer for Python telemetry service."""

import os
from enum import Enum
from typing import Protocol, runtime_checkable, Optional
from dataclasses import dataclass


class ProviderType(Enum):
    """Feature flag provider types."""
    UNLEASH = "unleash"
    LAUNCHDARKLY = "launchdarkly"


@runtime_checkable
class FeatureFlagProvider(Protocol):
    """Protocol for feature flag providers."""
    
    def is_enabled(self, flag_name: str, default_value: bool = False) -> bool:
        """Check if a feature flag is enabled."""
        ...
    
    def get_variant(self, flag_name: str, default_value: str = "") -> str:
        """Get the variant for a feature flag."""
        ...
    
    def close(self) -> None:
        """Close the provider connection."""
        ...


@dataclass
class UnleashConfig:
    """Unleash configuration."""
    url: str = None
    api_token: str = None
    app_name: str = None
    environment: str = None
    
    def __post_init__(self):
        if self.url is None:
            self.url = os.getenv("UNLEASH_URL", "http://localhost:4242")
        if self.api_token is None:
            self.api_token = os.getenv("UNLEASH_API_TOKEN", "")
        if self.app_name is None:
            self.app_name = os.getenv("UNLEASH_APP_NAME", "ci-cd-pipeline")
        if self.environment is None:
            self.environment = os.getenv("UNLEASH_ENVIRONMENT", "development")


@dataclass
class LaunchDarklyConfig:
    """LaunchDarkly configuration."""
    sdk_key: str = None
    app_name: str = None
    environment: str = None
    
    def __post_init__(self):
        if self.sdk_key is None:
            self.sdk_key = os.getenv("LAUNCHDARKLY_SDK_KEY", "")
        if self.app_name is None:
            self.app_name = os.getenv("LAUNCHDARKLY_APP_NAME", "ci-cd-pipeline")
        if self.environment is None:
            self.environment = os.getenv("LAUNCHDARKLY_ENVIRONMENT", "development")


@dataclass
class Config:
    """Feature flag configuration."""
    provider: ProviderType = None
    unleash: UnleashConfig = None
    launchdarkly: LaunchDarklyConfig = None
    
    def __post_init__(self):
        if self.provider is None:
            provider_str = os.getenv("FEATURE_FLAG_PROVIDER", "unleash")
            try:
                self.provider = ProviderType(provider_str)
            except ValueError:
                self.provider = ProviderType.UNLEASH
        
        if self.unleash is None:
            self.unleash = UnleashConfig()
        
        if self.launchdarkly is None:
            self.launchdarkly = LaunchDarklyConfig()


def create_provider(config: Optional[Config] = None) -> FeatureFlagProvider:
    """Create a feature flag provider based on configuration."""
    if config is None:
        config = Config()
    
    if config.provider == ProviderType.UNLEASH:
        from .unleash import UnleashProvider
        return UnleashProvider(config.unleash)
    elif config.provider == ProviderType.LAUNCHDARKLY:
        from .launchdarkly import LaunchDarklyProvider
        return LaunchDarklyProvider(config.launchdarkly)
    else:
        raise ValueError(f"Unknown provider type: {config.provider}")


# Feature flag names
FLAG_NEW_API_ENDPOINT = "new_api_endpoint"
FLAG_ADVANCED_TELEMETRY = "advanced_telemetry"
FLAG_EXPERIMENTAL_CACHE = "experimental_cache"
FLAG_RATE_LIMITING = "rate_limiting"
