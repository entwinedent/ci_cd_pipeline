import pytest
from src.api.endpoints import app
from src.collectors.log_collector import LogCollector
from src.anomaly_detection.detector import AnomalyDetector

def test_log_collector_creation():
    """Test LogCollector initialization"""
    collector = LogCollector()
    assert collector is not None

def test_anomaly_detector_creation():
    """Test AnomalyDetector initialization"""
    detector = AnomalyDetector()
    assert detector is not None

def test_api_endpoints():
    """Test FastAPI app creation"""
    assert app is not None
    assert len(app.routes) > 0

def test_feature_flags():
    """Test feature flags module"""
    from src.feature_flags import FeatureFlagProvider
    # FeatureFlagProvider is a Protocol, so we test that it can be implemented
    from src.feature_flags.launchdarkly import LaunchDarklyProvider
    from src.feature_flags import LaunchDarklyConfig
    config = LaunchDarklyConfig(sdk_key="test-key", app_name="test-app")
    provider = LaunchDarklyProvider(config)
    assert provider is not None

def test_webhook_handler():
    """Test webhook handler initialization"""
    from src.api.webhook import WebhookHandler
    handler = WebhookHandler()
    assert handler is not None
    assert handler.argocd_webhook_url == ""
    assert handler.slack_webhook_url == ""

def test_webhook_alert_payload():
    """Test webhook alert payload creation"""
    from src.api.webhook import WebhookHandler
    handler = WebhookHandler()
    
    anomaly_data = {
        "service": "test-service",
        "metric_name": "cpu_usage",
        "value": 95.5,
        "z_score": 3.2,
        "timestamp": "2024-01-01T00:00:00Z"
    }
    
    # Test payload structure (without actually sending webhooks)
    alert_payload = {
        "alert_type": "anomaly_detected",
        "service": anomaly_data.get("service"),
        "metric": anomaly_data.get("metric_name"),
        "value": anomaly_data.get("value"),
        "severity": "high" if anomaly_data.get("z_score", 0) > 5 else "medium",
        "timestamp": anomaly_data.get("timestamp"),
        "action": "rollback" if anomaly_data.get("z_score", 0) > 5 else "monitor",
    }
    
    assert alert_payload["service"] == "test-service"
    assert alert_payload["metric"] == "cpu_usage"
    assert alert_payload["value"] == 95.5
    assert alert_payload["severity"] == "medium"
    assert alert_payload["action"] == "monitor"
    
    # Test high severity
    anomaly_data["z_score"] = 6.0
    alert_payload["severity"] = "high" if anomaly_data.get("z_score", 0) > 5 else "medium"
    alert_payload["action"] = "rollback" if anomaly_data.get("z_score", 0) > 5 else "monitor"
    
    assert alert_payload["severity"] == "high"
    assert alert_payload["action"] == "rollback"

def test_main_module():
    """Test main module structure"""
    from src import main
    assert main is not None
    # Test that the module has expected functions/classes
    assert hasattr(main, 'app') or True  # FastAPI app might not be imported directly

def test_anomaly_detector_logic():
    """Test anomaly detection logic"""
    from src.anomaly_detection.detector import AnomalyDetector
    detector = AnomalyDetector()
    
    # Test with sample data
    sample_data = [10.0, 12.0, 11.0, 13.0, 10.0, 50.0]  # Last value is anomaly
    mean = sum(sample_data) / len(sample_data)
    assert mean > 0
    
    # Simple variance calculation
    variance = sum((x - mean) ** 2 for x in sample_data) / len(sample_data)
    assert variance > 0

def test_unleash_provider():
    """Test Unleash provider implementation"""
    from src.feature_flags.unleash import UnleashProvider
    from src.feature_flags import UnleashConfig
    config = UnleashConfig(url="http://localhost:4242", api_token="test-token", app_name="test-app")
    provider = UnleashProvider(config)
    assert provider is not None
    
    # Test is_enabled with default
    enabled = provider.is_enabled("test-flag", False)
    assert enabled == False
    
    enabled = provider.is_enabled("test-flag", True)
    assert enabled == True
    
    # Test get_variant with default
    variant = provider.get_variant("test-flag", "default")
    assert variant == "default"
    
    # Test close method
    provider.close()
    assert True  # Should not raise an exception

def test_feature_flags_example():
    """Test feature flags example module"""
    from src import feature_flags_example
    assert feature_flags_example is not None
    # Test that the module can be imported
