from fastapi import FastAPI, HTTPException, status
from pydantic import BaseModel
from typing import Optional, Dict, Any
import logging
from datetime import datetime

logger = logging.getLogger(__name__)

app = FastAPI(title="Telemetry Collector API", version="1.0.0")


class LogEntry(BaseModel):
    service: str
    level: str
    message: str
    timestamp: str
    metadata: Optional[Dict[str, Any]] = {}


class MetricEntry(BaseModel):
    service: str
    metric_name: str
    value: float
    timestamp: str
    labels: Optional[Dict[str, str]] = {}


class HealthResponse(BaseModel):
    status: str
    timestamp: str


class IngestResponse(BaseModel):
    success: bool
    message: str


# In-memory storage for demonstration
log_storage = []
metric_storage = []


@app.get("/healthz", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    return HealthResponse(
        status="healthy",
        timestamp=datetime.utcnow().isoformat()
    )


@app.get("/readyz", response_model=HealthResponse)
async def readiness_check():
    """Readiness check endpoint"""
    return HealthResponse(
        status="ready",
        timestamp=datetime.utcnow().isoformat()
    )


@app.post("/api/v1/logs", response_model=IngestResponse)
async def ingest_log(log_entry: LogEntry):
    """Ingest log entry"""
    try:
        log_storage.append(log_entry.dict())
        logger.info(f"Received log from {log_entry.service}: {log_entry.message}")
        
        # Trigger anomaly detection
        from ..anomaly_detection.detector import detect_anomalies
        anomalies = detect_anomalies(log_entry.service, log_entry.metadata)
        
        if anomalies:
            logger.warning(f"Anomalies detected for {log_entry.service}: {anomalies}")
            # Send webhook notifications
            await send_anomaly_alert(log_entry.service, anomalies)
        
        return IngestResponse(
            success=True,
            message="Log ingested successfully"
        )
    except Exception as e:
        logger.error(f"Failed to ingest log: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to ingest log: {str(e)}"
        )


@app.get("/api/v1/logs/{service}")
async def get_logs(service: str, limit: int = 100):
    """Retrieve logs for a specific service"""
    service_logs = [
        log for log in log_storage if log.get("service") == service
    ]
    return service_logs[-limit:]


@app.post("/api/v1/metrics", response_model=IngestResponse)
async def ingest_metric(metric_entry: MetricEntry):
    """Ingest metric entry"""
    try:
        metric_storage.append(metric_entry.dict())
        logger.info(f"Received metric from {metric_entry.service}: {metric_entry.metric_name} = {metric_entry.value}")
        
        return IngestResponse(
            success=True,
            message="Metric ingested successfully"
        )
    except Exception as e:
        logger.error(f"Failed to ingest metric: {e}")
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to ingest metric: {str(e)}"
        )


@app.get("/api/v1/metrics/{service}")
async def get_metrics(service: str, time_range: str = "1h"):
    """Retrieve metrics for a specific service"""
    service_metrics = [metric for metric in metric_storage if metric.get("service") == service]
    return service_metrics


async def send_anomaly_alert(service: str, anomalies: list):
    """Send anomaly alert via webhook"""
    import os
    import httpx
    
    webhook_url = os.getenv("SLACK_WEBHOOK_URL")
    if webhook_url:
        try:
            async with httpx.AsyncClient() as client:
                await client.post(
                    webhook_url,
                    json={
                        "text": f"Anomaly detected for service {service}",
                        "attachments": [
                            {
                                "text": f"Anomalies: {', '.join(anomalies)}",
                                "color": "#ff0000"
                            }
                        ]
                    }
                )
        except Exception as e:
            logger.error(f"Failed to send webhook alert: {e}")
