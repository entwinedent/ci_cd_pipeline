import logging
from typing import Dict, Any, List
from collections import deque
import statistics
from datetime import datetime, timedelta

logger = logging.getLogger(__name__)


class AnomalyDetector:
    def __init__(self, window_size: int = 100, threshold_std: float = 3.0):
        self.metric_history = {}  # service -> metric_name -> deque of values
        self.window_size = window_size
        self.threshold_std = threshold_std
        self.detected_anomalies = deque(maxlen=1000)
    
    def detect_anomaly(self, metric_data: Dict[str, Any]) -> bool:
        """Detect if a metric value is anomalous using statistical analysis"""
        service = metric_data.get('service', 'unknown')
        metric_name = metric_data.get('metric_name', 'unknown')
        value = metric_data.get('value', 0.0)
        
        # Initialize history for this metric if needed
        key = f"{service}:{metric_name}"
        if key not in self.metric_history:
            self.metric_history[key] = deque(maxlen=self.window_size)
        
        history = self.metric_history[key]
        
        # Need at least some data points for anomaly detection
        if len(history) < 10:
            history.append(value)
            return False
        
        # Calculate statistics using Python's statistics module
        values = list(history)
        try:
            mean = statistics.mean(values)
            std = statistics.stdev(values)
        except statistics.StatisticsError:
            # Not enough variance or data points
            history.append(value)
            return False
        
        # Z-score based anomaly detection
        if std > 0:
            z_score = abs(value - mean) / std
            is_anomaly = z_score > self.threshold_std
        else:
            is_anomaly = False
        
        # Add to history
        history.append(value)
        
        # Record anomaly if detected
        if is_anomaly:
            anomaly_record = {
                'service': service,
                'metric_name': metric_name,
                'value': value,
                'mean': mean,
                'std': std,
                'z_score': z_score if std > 0 else 0,
                'timestamp': datetime.utcnow().isoformat()
            }
            self.detected_anomalies.append(anomaly_record)
            logger.warning(f"Anomaly detected: {anomaly_record}")
        
        return is_anomaly
    
    def get_recent_anomalies(self, limit: int = 100) -> List[Dict[str, Any]]:
        """Get recently detected anomalies"""
        return list(self.detected_anomalies)[-limit:]
    
    def get_anomalies_by_service(self, service: str, limit: int = 100) -> List[Dict[str, Any]]:
        """Get anomalies for a specific service"""
        service_anomalies = [
            anomaly for anomaly in self.detected_anomalies 
            if anomaly.get('service') == service
        ]
        return service_anomalies[-limit:]
