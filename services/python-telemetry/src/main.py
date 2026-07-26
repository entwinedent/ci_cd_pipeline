import os
import logging
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import Optional, Dict, Any
import uvicorn
from src.collectors.log_collector import LogCollector
from src.anomaly_detection.detector import AnomalyDetector
from src.api.webhook import WebhookHandler

# Configure logging
logging.basicConfig(
    level=os.getenv("LOG_LEVEL", "INFO").upper(),
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = FastAPI(title="Telemetry Collector", version="1.0.0")

# Initialize components
log_collector = LogCollector()
anomaly_detector = AnomalyDetector()
webhook_handler = WebhookHandler()


class LogEntry(BaseModel):
    service: str
    level: str
    message: str
    timestamp: str
    metadata: Optional[Dict[str, Any]] = None


class MetricData(BaseModel):
    service: str
    metric_name: str
    value: float
    timestamp: str
    labels: Optional[Dict[str, str]] = None


class HealthResponse(BaseModel):
    status: str
    message: str


@app.get("/healthz")
async def healthz() -> HealthResponse:
    """Liveness probe endpoint"""
    return HealthResponse(status="healthy", message="Service is running")


@app.get("/readyz")
async def readyz() -> HealthResponse:
    """Readiness probe endpoint"""
    # TODO: Check connections to other services
    return HealthResponse(status="ready", message="Service is ready to accept traffic")


@app.post("/api/v1/logs")
async def ingest_logs(log_entry: LogEntry):
    """Ingest structured logs from services"""
    try:
        log_collector.process_log(log_entry.dict())
        logger.info(f"Processed log from {log_entry.service}")
        return {"status": "success", "message": "Log ingested successfully"}
    except Exception as e:
        logger.error(f"Failed to process log: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/metrics")
async def ingest_metrics(metric_data: MetricData):
    """Ingest metrics from services and check for anomalies"""
    try:
        is_anomaly = anomaly_detector.detect_anomaly(metric_data.dict())

        if is_anomaly:
            logger.warning(
                f"Anomaly detected in {metric_data.service}: "
                f"{metric_data.metric_name}"
            )
            # Trigger webhook for auto-remediation
            await webhook_handler.trigger_alert(metric_data.dict())
        
        return {
            "status": "success",
            "anomaly_detected": is_anomaly,
            "message": "Metric processed successfully"
        }
    except Exception as e:
        logger.error(f"Failed to process metric: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/api/v1/anomalies")
async def get_recent_anomalies():
    """Get recent detected anomalies"""
    try:
        anomalies = anomaly_detector.get_recent_anomalies()
        return {"anomalies": anomalies}
    except Exception as e:
        logger.error(f"Failed to retrieve anomalies: {e}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    port = int(os.getenv("PORT", "8000"))
    uvicorn.run(app, host="0.0.0.0", port=port)
