use std::env;

/// Feature flag provider types
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ProviderType {
    Unleash,
    LaunchDarkly,
}

impl std::str::FromStr for ProviderType {
    type Err = String;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        match s.to_lowercase().as_str() {
            "unleash" => Ok(ProviderType::Unleash),
            "launchdarkly" => Ok(ProviderType::LaunchDarkly),
            _ => Err(format!("Unknown provider type: {}", s)),
        }
    }
}

/// Feature flag provider trait
pub trait FeatureFlagProvider: Send + Sync {
    /// Check if a feature flag is enabled
    fn is_enabled(&self, flag_name: &str, default_value: bool) -> Result<bool, Box<dyn std::error::Error>>;
    
    /// Get the variant for a feature flag
    fn get_variant(&self, flag_name: &str, default_value: &str) -> Result<String, Box<dyn std::error::Error>>;
}

/// Configuration for feature flag providers
#[derive(Debug, Clone)]
pub struct Config {
    pub provider: ProviderType,
    pub unleash: UnleashConfig,
    pub launchdarkly: LaunchDarklyConfig,
}

impl Default for Config {
    fn default() -> Self {
        let provider_env = env::var("FEATURE_FLAG_PROVIDER")
            .unwrap_or_else(|_| "unleash".to_string());
        
        Config {
            provider: provider_env.parse().unwrap_or(ProviderType::Unleash),
            unleash: UnleashConfig::default(),
            launchdarkly: LaunchDarklyConfig::default(),
        }
    }
}

/// Unleash configuration
#[derive(Debug, Clone)]
pub struct UnleashConfig {
    pub url: String,
    pub api_token: String,
    pub app_name: String,
    pub environment: String,
}

impl Default for UnleashConfig {
    fn default() -> Self {
        UnleashConfig {
            url: env::var("UNLEASH_URL").unwrap_or_else(|_| "http://localhost:4242".to_string()),
            api_token: env::var("UNLEASH_API_TOKEN").unwrap_or_else(|_| "".to_string()),
            app_name: env::var("UNLEASH_APP_NAME").unwrap_or_else(|_| "ci-cd-pipeline".to_string()),
            environment: env::var("UNLEASH_ENVIRONMENT").unwrap_or_else(|_| "development".to_string()),
        }
    }
}

/// LaunchDarkly configuration
#[derive(Debug, Clone)]
pub struct LaunchDarklyConfig {
    pub sdk_key: String,
    pub app_name: String,
    pub environment: String,
}

impl Default for LaunchDarklyConfig {
    fn default() -> Self {
        LaunchDarklyConfig {
            sdk_key: env::var("LAUNCHDARKLY_SDK_KEY").unwrap_or_else(|_| "".to_string()),
            app_name: env::var("LAUNCHDARKLY_APP_NAME").unwrap_or_else(|_| "ci-cd-pipeline".to_string()),
            environment: env::var("LAUNCHDARKLY_ENVIRONMENT").unwrap_or_else(|_| "development".to_string()),
        }
    }
}

/// Create a new feature flag provider based on configuration
pub fn create_provider(config: Config) -> Result<Box<dyn FeatureFlagProvider>, String> {
    match config.provider {
        ProviderType::Unleash => Ok(Box::new(UnleashProvider::new(config.unleash))),
        ProviderType::LaunchDarkly => Ok(Box::new(LaunchDarklyProvider::new(config.launchdarkly)?)),
    }
}

/// Feature flag names
pub const FLAG_NEW_API_ENDPOINT: &str = "new_api_endpoint";
pub const FLAG_ADVANCED_TELEMETRY: &str = "advanced_telemetry";
pub const FLAG_EXPERIMENTAL_CACHE: &str = "experimental_cache";
pub const FLAG_RATE_LIMITING: &str = "rate_limiting";

pub mod unleash;
pub mod launchdarkly;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_launchdarkly_provider_structure() {
        let config = LaunchDarklyConfig {
            sdk_key: "test-sdk-key".to_string(),
            app_name: "test-app".to_string(),
            environment: "development".to_string(),
        };
        let provider = launchdarkly::LaunchDarklyProvider::new(config);
        assert!(provider.is_ok());
    }

    #[test]
    fn test_launchdarkly_provider_validation() {
        let config = LaunchDarklyConfig {
            sdk_key: "".to_string(),
            app_name: "test-app".to_string(),
            environment: "development".to_string(),
        };
        let provider = launchdarkly::LaunchDarklyProvider::new(config);
        assert!(provider.is_err());
    }

    #[test]
    fn test_config_default() {
        let config = Config::default();
        assert!(matches!(config.provider, ProviderType::Unleash));
    }

    #[test]
    fn test_provider_type_parsing() {
        assert_eq!("unleash".parse::<ProviderType>(), Ok(ProviderType::Unleash));
        assert_eq!("launchdarkly".parse::<ProviderType>(), Ok(ProviderType::LaunchDarkly));
        assert!("unknown".parse::<ProviderType>().is_err());
    }
}
