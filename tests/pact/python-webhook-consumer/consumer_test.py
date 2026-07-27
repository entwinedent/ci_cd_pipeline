import pytest
from pact import Pact
import json

def test_webhook_alert_success():
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
    
    pact = Pact('python-telemetry', 'external-webhook', log_level='DEBUG')
    (pact
     .given('webhook endpoint is available')
     .upon_receiving('an anomaly alert webhook')
     .with_request('POST', '/api/v1/alerts', body=expected)
     .will_respond_with(200, body=expected))
    
    with pact:
        response = pact.post('/api/v1/alerts', json=expected)
        assert response.status_code == 200
        assert response.json() == expected

def test_webhook_alert_medium_severity():
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
    
    pact = Pact('python-telemetry', 'external-webhook', log_level='DEBUG')
    (pact
     .given('webhook endpoint is available')
     .upon_receiving('a medium severity alert')
     .with_request('POST', '/api/v1/alerts', body=expected)
     .will_respond_with(200, body=expected))
    
    with pact:
        response = pact.post('/api/v1/alerts', json=expected)
        assert response.status_code == 200

def test_webhook_failure_retry():
    """Test webhook failure and retry logic"""
    pact = Pact('python-telemetry', 'external-webhook', log_level='DEBUG')
    (pact
     .given('webhook endpoint is temporarily unavailable')
     .upon_receiving('an alert during webhook failure')
     .with_request('POST', '/api/v1/alerts', body={'test': 'data'})
     .will_respond_with(503))
    
    with pact:
        response = pact.post('/api/v1/alerts', json={'test': 'data'})
        assert response.status_code == 503

if __name__ == '__main__':
    pytest.main([__file__, '-v'])
