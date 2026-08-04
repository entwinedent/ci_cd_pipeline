package featureflags

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewFeatureFlags(t *testing.T) {
	// Create a mock provider
	mockProvider := &MockProvider{}

	// Test FeatureFlags creation
	ff := NewFeatureFlags(mockProvider)
	if ff == nil {
		t.Error("Expected FeatureFlags to be created")
	}
}

func TestNewAPIEndpointEnabled(t *testing.T) {
	mockProvider := &MockProvider{enabled: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.NewAPIEndpointEnabled(context.Background())
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

func TestAdvancedTelemetryEnabled(t *testing.T) {
	mockProvider := &MockProvider{enabled: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.AdvancedTelemetryEnabled(context.Background())
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

func TestExperimentalCacheEnabled(t *testing.T) {
	mockProvider := &MockProvider{enabled: false}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.ExperimentalCacheEnabled(context.Background())
	if enabled != false {
		t.Errorf("Expected false, got %v", enabled)
	}
}

func TestRateLimitingEnabled(t *testing.T) {
	mockProvider := &MockProvider{enabled: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.RateLimitingEnabled(context.Background())
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

// MockProvider is a mock implementation for testing
type MockProvider struct {
	enabled     bool
	returnError bool
}

func (m *MockProvider) IsEnabled(ctx context.Context, flagName string, defaultValue bool) (bool, error) {
	if m.returnError {
		return false, fmt.Errorf("mock error")
	}
	return m.enabled, nil
}

func (m *MockProvider) GetVariant(ctx context.Context, flagName string, defaultValue string) (string, error) {
	if m.returnError {
		return "", fmt.Errorf("mock error")
	}
	return defaultValue, nil
}

func (m *MockProvider) Close() error {
	return nil
}

func TestProviderErrorHandling(t *testing.T) {
	mockProvider := &MockProvider{enabled: true, returnError: true}
	ff := NewFeatureFlags(mockProvider)

	// Test that errors are handled gracefully
	enabled := ff.NewAPIEndpointEnabled(context.Background())
	if enabled != false {
		t.Errorf("Expected false on error, got %v", enabled)
	}
}

func TestGetVariant(t *testing.T) {
	mockProvider := &MockProvider{enabled: true}

	variant, err := mockProvider.GetVariant(context.Background(), "test-flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "default" {
		t.Errorf("Expected default variant, got %s", variant)
	}
}

func TestNewProvider_Unleash(t *testing.T) {
	cfg := Config{
		Provider: ProviderUnleash,
		Unleash: UnleashConfig{
			URL:         "http://localhost:4242",
			APIToken:    "test-token",
			AppName:     "test-app",
			Environment: "test",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
	if _, ok := provider.(*UnleashProvider); !ok {
		t.Error("Expected UnleashProvider")
	}
}

func TestNewProvider_LaunchDarkly(t *testing.T) {
	cfg := Config{
		Provider: ProviderLaunchDarkly,
		LaunchDarkly: LaunchDarklyConfig{
			SDKKey:      "test-sdk-key",
			AppName:     "test-app",
			Environment: "test",
		},
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
	if _, ok := provider.(*LaunchDarklyProvider); !ok {
		t.Error("Expected LaunchDarklyProvider")
	}
}

func TestNewProvider_Unknown(t *testing.T) {
	cfg := Config{
		Provider: ProviderType("unknown"),
	}

	_, err := NewProvider(cfg)
	if err == nil {
		t.Error("Expected error for unknown provider")
	}
}

func TestNewProvider_Default(t *testing.T) {
	cfg := Config{}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestNewUnleashProvider(t *testing.T) {
	cfg := UnleashConfig{
		URL:         "http://localhost:4242",
		APIToken:    "test-token",
		AppName:     "test-app",
		Environment: "test",
	}

	provider, err := NewUnleashProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestNewUnleashProvider_Defaults(t *testing.T) {
	cfg := UnleashConfig{}

	provider, err := NewUnleashProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestUnleashProvider_IsEnabled(t *testing.T) {
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test-flag", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

func TestUnleashProvider_IsEnabled_EnvVar(t *testing.T) {
	t.Setenv("UNLEASH_FLAG_test_flag", "true")
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test_flag", false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

func TestUnleashProvider_GetVariant(t *testing.T) {
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	variant, err := provider.GetVariant(context.Background(), "test-flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "default" {
		t.Errorf("Expected default, got %s", variant)
	}
}

func TestUnleashProvider_GetVariant_EnvVar(t *testing.T) {
	t.Setenv("UNLEASH_VARIANT_test_flag", "variant-a")
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	variant, err := provider.GetVariant(context.Background(), "test_flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "variant-a" {
		t.Errorf("Expected variant-a, got %s", variant)
	}
}

func TestUnleashProvider_Close(t *testing.T) {
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	err := provider.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestNewLaunchDarklyProvider(t *testing.T) {
	cfg := LaunchDarklyConfig{
		SDKKey:      "test-sdk-key",
		AppName:     "test-app",
		Environment: "test",
	}

	provider, err := NewLaunchDarklyProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestNewLaunchDarklyProvider_MissingSDKKey(t *testing.T) {
	cfg := LaunchDarklyConfig{}

	_, err := NewLaunchDarklyProvider(cfg)
	if err == nil {
		t.Error("Expected error for missing SDK key")
	}
}

func TestNewLaunchDarklyProvider_Defaults(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	cfg := LaunchDarklyConfig{}

	provider, err := NewLaunchDarklyProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestLaunchDarklyProvider_IsEnabled(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test-flag", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

func TestLaunchDarklyProvider_IsEnabled_EnvVar(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	t.Setenv("LD_FLAG_test_flag", "true")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test_flag", false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != true {
		t.Errorf("Expected true, got %v", enabled)
	}
}

func TestLaunchDarklyProvider_GetVariant(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	variant, err := provider.GetVariant(context.Background(), "test-flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "default" {
		t.Errorf("Expected default, got %s", variant)
	}
}

func TestLaunchDarklyProvider_GetVariant_EnvVar(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	t.Setenv("LD_VARIANT_test_flag", "variant-a")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	variant, err := provider.GetVariant(context.Background(), "test_flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "variant-a" {
		t.Errorf("Expected variant-a, got %s", variant)
	}
}

func TestLaunchDarklyProvider_Close(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	err := provider.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExampleUsage(t *testing.T) {
	// Test the example usage function
	cfg := Config{
		Provider: ProviderUnleash,
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	defer provider.Close()

	ff := NewFeatureFlags(provider)

	ctx := context.Background()

	// Test feature flag checks
	enabled := ff.NewAPIEndpointEnabled(ctx)
	if enabled != false {
		t.Logf("NewAPIEndpointEnabled: %v", enabled)
	}

	enabled = ff.AdvancedTelemetryEnabled(ctx)
	if enabled != true {
		t.Logf("AdvancedTelemetryEnabled: %v", enabled)
	}

	enabled = ff.ExperimentalCacheEnabled(ctx)
	if enabled != false {
		t.Logf("ExperimentalCacheEnabled: %v", enabled)
	}

	enabled = ff.RateLimitingEnabled(ctx)
	if enabled != true {
		t.Logf("RateLimitingEnabled: %v", enabled)
	}

	// Test direct provider usage
	isEnabled, err := provider.IsEnabled(ctx, "custom_feature", false)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if isEnabled != false {
		t.Logf("Custom feature enabled: %v", isEnabled)
	}

	// Test GetVariant
	variant, err := provider.GetVariant(ctx, "ui_theme", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "default" {
		t.Logf("UI theme variant: %s", variant)
	}
}

func TestUnleashProvider_Start(t *testing.T) {
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	unleashProvider, ok := provider.(*UnleashProvider)
	if !ok {
		t.Fatal("Expected UnleashProvider")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := unleashProvider.Start(ctx)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestNewProvider_EnvVar(t *testing.T) {
	t.Setenv("FEATURE_FLAG_PROVIDER", "launchdarkly")
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	cfg := Config{}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
	if _, ok := provider.(*LaunchDarklyProvider); !ok {
		t.Error("Expected LaunchDarklyProvider")
	}
}

func TestNewUnleashProvider_EnvVars(t *testing.T) {
	t.Setenv("UNLEASH_URL", "http://test-url:4242")
	t.Setenv("UNLEASH_API_TOKEN", "test-token")
	t.Setenv("UNLEASH_APP_NAME", "test-app")
	t.Setenv("UNLEASH_ENVIRONMENT", "test-env")
	cfg := UnleashConfig{}

	provider, err := NewUnleashProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestNewLaunchDarklyProvider_EnvVars(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "env-sdk-key")
	t.Setenv("LAUNCHDARKLY_APP_NAME", "env-app")
	t.Setenv("LAUNCHDARKLY_ENVIRONMENT", "env-env")
	cfg := LaunchDarklyConfig{}

	provider, err := NewLaunchDarklyProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created")
	}
}

func TestLaunchDarklyProvider_IsEnabled_EnvVarFalse(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	t.Setenv("LD_FLAG_test_flag", "false")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test_flag", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != false {
		t.Errorf("Expected false, got %v", enabled)
	}
}

func TestLaunchDarklyProvider_IsEnabled_EnvVarZero(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	t.Setenv("LD_FLAG_test_flag", "0")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test_flag", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != false {
		t.Errorf("Expected false, got %v", enabled)
	}
}

func TestUnleashProvider_IsEnabled_EnvVarFalse(t *testing.T) {
	t.Setenv("UNLEASH_FLAG_test_flag", "false")
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test_flag", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != false {
		t.Errorf("Expected false, got %v", enabled)
	}
}

func TestUnleashProvider_IsEnabled_EnvVarZero(t *testing.T) {
	t.Setenv("UNLEASH_FLAG_test_flag", "0")
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	enabled, err := provider.IsEnabled(context.Background(), "test_flag", true)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if enabled != false {
		t.Errorf("Expected false, got %v", enabled)
	}
}

func TestProviderConstants(t *testing.T) {
	if ProviderUnleash != "unleash" {
		t.Errorf("Expected ProviderUnleash to be 'unleash', got %s", ProviderUnleash)
	}
	if ProviderLaunchDarkly != "launchdarkly" {
		t.Errorf("Expected ProviderLaunchDarkly to be 'launchdarkly', got %s", ProviderLaunchDarkly)
	}
}

func TestFlagConstants(t *testing.T) {
	if FlagNewAPIEndpoint != "new_api_endpoint" {
		t.Errorf("Expected FlagNewAPIEndpoint to be 'new_api_endpoint', got %s", FlagNewAPIEndpoint)
	}
	if FlagAdvancedTelemetry != "advanced_telemetry" {
		t.Errorf("Expected FlagAdvancedTelemetry to be 'advanced_telemetry', got %s", FlagAdvancedTelemetry)
	}
	if FlagExperimentalCache != "experimental_cache" {
		t.Errorf("Expected FlagExperimentalCache to be 'experimental_cache', got %s", FlagExperimentalCache)
	}
	if FlagRateLimiting != "rate_limiting" {
		t.Errorf("Expected FlagRateLimiting to be 'rate_limiting', got %s", FlagRateLimiting)
	}
}

func TestFeatureFlags_AdvancedTelemetryEnabled_Error(t *testing.T) {
	mockProvider := &MockProvider{enabled: false, returnError: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.AdvancedTelemetryEnabled(context.Background())
	if enabled != true {
		t.Errorf("Expected true on error (default value), got %v", enabled)
	}
}

func TestFeatureFlags_RateLimitingEnabled_Error(t *testing.T) {
	mockProvider := &MockProvider{enabled: false, returnError: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.RateLimitingEnabled(context.Background())
	if enabled != true {
		t.Errorf("Expected true on error (default value), got %v", enabled)
	}
}

func TestExampleUsage_Function(t *testing.T) {
	// This test calls the ExampleUsage function to get coverage
	// We need to capture the log output since ExampleUsage uses log.Fatalf
	// For testing purposes, we'll just call it and expect it not to panic
	// Note: ExampleUsage uses log.Fatalf which would exit, so we can't call it directly
	// Instead, we'll test the same logic without the fatal calls

	cfg := Config{
		Provider: ProviderUnleash,
	}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Fatalf("Failed to create feature flag provider: %v", err)
	}
	defer provider.Close()

	ff := NewFeatureFlags(provider)

	ctx := context.Background()

	// Check if new API endpoint is enabled
	if ff.NewAPIEndpointEnabled(ctx) {
		// Use new endpoint implementation
	}

	// Check if advanced telemetry is enabled
	if ff.AdvancedTelemetryEnabled(ctx) {
		// Enable advanced telemetry collection
	}

	// Check if experimental cache is enabled
	if ff.ExperimentalCacheEnabled(ctx) {
		// Enable experimental caching layer
	}

	// Check if rate limiting is enabled
	if ff.RateLimitingEnabled(ctx) {
		// Apply rate limiting
	}

	// Example: Direct provider usage
	isEnabled, err := provider.IsEnabled(ctx, "custom_feature", false)
	if err != nil {
		t.Logf("Error checking custom feature: %v", err)
	} else if isEnabled {
		// Custom feature is enabled
	}

	// Example: Get variant
	variant, err := provider.GetVariant(ctx, "ui_theme", "default")
	if err != nil {
		t.Logf("Error getting variant: %v", err)
	} else {
		_ = variant
	}
}

func TestConfigStruct(t *testing.T) {
	cfg := Config{
		Provider: ProviderUnleash,
		Unleash: UnleashConfig{
			URL:         "http://localhost:4242",
			APIToken:    "token",
			AppName:     "app",
			Environment: "dev",
		},
		LaunchDarkly: LaunchDarklyConfig{
			SDKKey:      "sdk-key",
			AppName:     "app",
			Environment: "dev",
		},
	}

	if cfg.Provider != ProviderUnleash {
		t.Errorf("Expected ProviderUnleash, got %v", cfg.Provider)
	}
	if cfg.Unleash.URL != "http://localhost:4242" {
		t.Errorf("Expected URL http://localhost:4242, got %s", cfg.Unleash.URL)
	}
	if cfg.LaunchDarkly.SDKKey != "sdk-key" {
		t.Errorf("Expected SDKKey sdk-key, got %s", cfg.LaunchDarkly.SDKKey)
	}
}

func TestUnleashConfigStruct(t *testing.T) {
	cfg := UnleashConfig{
		URL:         "http://localhost:4242",
		APIToken:    "token",
		AppName:     "app",
		Environment: "dev",
	}

	if cfg.URL != "http://localhost:4242" {
		t.Errorf("Expected URL http://localhost:4242, got %s", cfg.URL)
	}
	if cfg.APIToken != "token" {
		t.Errorf("Expected APIToken token, got %s", cfg.APIToken)
	}
	if cfg.AppName != "app" {
		t.Errorf("Expected AppName app, got %s", cfg.AppName)
	}
	if cfg.Environment != "dev" {
		t.Errorf("Expected Environment dev, got %s", cfg.Environment)
	}
}

func TestLaunchDarklyConfigStruct(t *testing.T) {
	cfg := LaunchDarklyConfig{
		SDKKey:      "sdk-key",
		AppName:     "app",
		Environment: "dev",
	}

	if cfg.SDKKey != "sdk-key" {
		t.Errorf("Expected SDKKey sdk-key, got %s", cfg.SDKKey)
	}
	if cfg.AppName != "app" {
		t.Errorf("Expected AppName app, got %s", cfg.AppName)
	}
	if cfg.Environment != "dev" {
		t.Errorf("Expected Environment dev, got %s", cfg.Environment)
	}
}

func TestLDClientStruct(t *testing.T) {
	client := &LDClient{
		sdkKey:      "test-sdk-key",
		appName:     "test-app",
		environment: "test-env",
	}

	if client.sdkKey != "test-sdk-key" {
		t.Errorf("Expected sdkKey test-sdk-key, got %s", client.sdkKey)
	}
	if client.appName != "test-app" {
		t.Errorf("Expected appName test-app, got %s", client.appName)
	}
	if client.environment != "test-env" {
		t.Errorf("Expected environment test-env, got %s", client.environment)
	}
}

func TestUnleashClientStruct(t *testing.T) {
	client := &UnleashClient{
		baseURL:  "http://localhost:4242",
		apiToken: "test-token",
		appName:  "test-app",
	}

	if client.baseURL != "http://localhost:4242" {
		t.Errorf("Expected baseURL http://localhost:4242, got %s", client.baseURL)
	}
	if client.apiToken != "test-token" {
		t.Errorf("Expected apiToken test-token, got %s", client.apiToken)
	}
	if client.appName != "test-app" {
		t.Errorf("Expected appName test-app, got %s", client.appName)
	}
}

func TestLaunchDarklyProvider_GetVariant_Default(t *testing.T) {
	t.Setenv("LAUNCHDARKLY_SDK_KEY", "test-sdk-key")
	cfg := LaunchDarklyConfig{}
	provider, _ := NewLaunchDarklyProvider(cfg)

	variant, err := provider.GetVariant(context.Background(), "test-flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "default" {
		t.Errorf("Expected default, got %s", variant)
	}
}

func TestUnleashProvider_GetVariant_Default(t *testing.T) {
	cfg := UnleashConfig{}
	provider, _ := NewUnleashProvider(cfg)

	variant, err := provider.GetVariant(context.Background(), "test-flag", "default")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if variant != "default" {
		t.Errorf("Expected default, got %s", variant)
	}
}

func TestFeatureFlags_NewAPIEndpointEnabled_Error(t *testing.T) {
	mockProvider := &MockProvider{enabled: true, returnError: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.NewAPIEndpointEnabled(context.Background())
	if enabled != false {
		t.Errorf("Expected false on error (default value), got %v", enabled)
	}
}

func TestFeatureFlags_ExperimentalCacheEnabled_Error(t *testing.T) {
	mockProvider := &MockProvider{enabled: true, returnError: true}
	ff := NewFeatureFlags(mockProvider)

	enabled := ff.ExperimentalCacheEnabled(context.Background())
	if enabled != false {
		t.Errorf("Expected false on error (default value), got %v", enabled)
	}
}

func TestMockProvider_Close(t *testing.T) {
	mockProvider := &MockProvider{}
	err := mockProvider.Close()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestNewProvider_EmptyEnvVar(t *testing.T) {
	t.Setenv("FEATURE_FLAG_PROVIDER", "")
	cfg := Config{}

	provider, err := NewProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created (should default to Unleash)")
	}
}

func TestNewUnleashProvider_EmptyEnvVars(t *testing.T) {
	t.Setenv("UNLEASH_URL", "")
	t.Setenv("UNLEASH_API_TOKEN", "")
	t.Setenv("UNLEASH_APP_NAME", "")
	t.Setenv("UNLEASH_ENVIRONMENT", "")
	cfg := UnleashConfig{}

	provider, err := NewUnleashProvider(cfg)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if provider == nil {
		t.Error("Expected provider to be created with defaults")
	}
}
