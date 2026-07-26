import logging
import os
import httpx
from typing import Dict, Any

logger = logging.getLogger(__name__)


class WebhookHandler:
    def __init__(self):
        self.argocd_webhook_url = os.getenv("ARGOCD_WEBHOOK_URL", "")
        self.slack_webhook_url = os.getenv("SLACK_WEBHOOK_URL", "")
    
    async def trigger_alert(self, anomaly_data: Dict[str, Any]) -> bool:
        """Trigger alert for detected anomaly"""
        try:
            # Prepare alert payload
            alert_payload = {
                'alert_type': 'anomaly_detected',
                'service': anomaly_data.get('service'),
                'metric': anomaly_data.get('metric_name'),
                'value': anomaly_data.get('value'),
                'severity': 'high' if anomaly_data.get('z_score', 0) > 5 else 'medium',
                'timestamp': anomaly_data.get('timestamp'),
                'action': 'rollback' if anomaly_data.get('z_score', 0) > 5 else 'monitor'
            }
            
            # Send to Argo CD if configured (for auto-remediation)
            if self.argocd_webhook_url:
                await self._send_webhook(self.argocd_webhook_url, alert_payload)
            
            # Send to Slack if configured (for notification)
            if self.slack_webhook_url:
                await self._send_slack_alert(alert_payload)
            
            logger.info(f"Alert triggered for {anomaly_data.get('service')}")
            return True
            
        except Exception as e:
            logger.error(f"Failed to trigger alert: {e}")
            return False
    
    async def _send_webhook(self, url: str, payload: Dict[str, Any]) -> None:
        """Send webhook to specified URL"""
        async with httpx.AsyncClient() as client:
            response = await client.post(url, json=payload, timeout=10.0)
            response.raise_for_status()
    
    async def _send_slack_alert(self, payload: Dict[str, Any]) -> None:
        """Send formatted alert to Slack"""
        slack_message = {
            'text': f"🚨 Anomaly Detected in {payload['service']}",
            'attachments': [{
                'color': 'danger' if payload['severity'] == 'high' else 'warning',
                'fields': [
                    {'title': 'Service', 'value': payload['service'], 'short': True},
                    {'title': 'Metric', 'value': payload['metric'], 'short': True},
                    {'title': 'Value', 'value': str(payload['value']), 'short': True},
                    {'title': 'Severity', 'value': payload['severity'], 'short': True},
                    {'title': 'Action', 'value': payload['action'], 'short': True},
                ]
            }]
        }
        
        await self._send_webhook(self.slack_webhook_url, slack_message)
