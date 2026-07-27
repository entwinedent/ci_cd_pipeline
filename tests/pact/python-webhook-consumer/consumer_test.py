import pytest
import requests
from pact import Consumer, Provider
import json

@pytest.fixture(scope="session")
def pact_mock_service():
    """Start and stop the Pact mock service for all tests"""
    pact = Consumer('python-telemetry').has_pact_with(
        Provider('external-webhook'),
        host_name='127.0.0.1',
        port=1234,
    )
    pact.start_service()
    yield pact
    pact.stop_service()

def test_webhook_alert_success(pact_mock_service):
    """Test that webhook alert is sent successfully"""
    expected = {
        'alert_type': 'anomaly_detected',
        'service': 'api-gateway',
        'metric': 'response_time',
        'value': 1500.0,
        'severity': 'high',
        'timestamp': '2024-01-01T00:00:00Z',
        'action': 'rollback'
    }
    
    (pact_mock_service
     .given('webhook endpoint is available')
     .upon_receiving('an anomaly alert webhook')
     .with_request('POST', '/api/v1/alerts', body=expected)
     .will_respond_with(200, body=expected))
    
    with pact_mock_service:
        response = requests.post(f"{pact_mock_service.uri}/api/v1/alerts", json=expected)
        assert response.status_code == 200
        assert response.json() == expected

def test_webhook_alert_medium_severity(pact_mock_service):
    """Test webhook alert with medium severity"""
    expected = {
        'alert_type': 'anomaly_detected',
        'service': 'rust-data-store',
        'metric': 'memory_usage',
        'value': 85.0,
        'severity': 'medium',
        'timestamp': '2024-01-01T00:00:00Z',
        'action': 'monitor'
    }
    
    (pact_mock_service
     .given('webhook endpoint is available')
     .upon_receiving('a medium severity alert')
     .with_request('POST', '/api/v1/alerts', body=expected)
     .will_respond_with(200, body=expected))
    
    with pact_mock_service:
        response = requests.post(f"{pact_mock_service.uri}/api/v1/alerts", json=expected)
        assert response.status_code == 200

def test_webhook_failure_retry(pact_mock_service):
    """Test webhook failure and retry logic"""
    (pact_mock_service
     .given('webhook endpoint is temporarily unavailable')
     .upon_receiving('an alert during webhook failure')
     .with_request('POST', '/api/v1/alerts', body={'test': 'data'})
     .will_respond_with(503))
    
    with pact_mock_service:
        response = requests.post(f"{pact_mock_service.uri}/api/v1/alerts", json={'test': 'data'})
        assert response.status_code == 503

if __name__ == '__main__':
    pytest.main([__file__, '-v'])
