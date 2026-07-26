import logging
from typing import Dict, Any
from collections import deque

logger = logging.getLogger(__name__)


class LogCollector:
    def __init__(self, max_entries: int = 10000):
        self.log_buffer: deque[Dict[str, Any]] = deque(maxlen=max_entries)
        self.max_entries = max_entries
    
    def process_log(self, log_entry: Dict[str, Any]) -> None:
        """Process and store log entry"""
        try:
            # Add timestamp if not present
            if 'timestamp' not in log_entry:
                from datetime import datetime
                log_entry['timestamp'] = datetime.utcnow().isoformat()
            
            # Store in buffer
            self.log_buffer.append(log_entry)
            logger.debug(f"Collected log from {log_entry.get('service', 'unknown')}")
        except Exception as e:
            logger.error(f"Failed to process log entry: {e}")
    
    def get_recent_logs(self, limit: int = 100) -> list:
        """Get recent log entries"""
        return list(self.log_buffer)[-limit:]
    
    def get_logs_by_service(self, service: str, limit: int = 100) -> list:
        """Get logs for a specific service"""
        service_logs = [
            log for log in self.log_buffer
            if log.get('service') == service
        ]
        return service_logs[-limit:]
    
    def get_error_logs(self, limit: int = 100) -> list:
        """Get error-level logs"""
        error_logs = [
            log for log in self.log_buffer
            if log.get('level') == 'ERROR'
        ]
        return error_logs[-limit:]
