"""
Integration tests for the CI/CD pipeline services
Tests end-to-end communication between Go API Gateway, Rust Data Store, and Python Telemetry
"""

import pytest
import grpc
import requests
import time
from typing import Dict, Any
import data_store_pb2
import data_store_pb2_grpc


class TestGoAPIGateway:
    """Integration tests for Go API Gateway"""
    
    BASE_URL = "http://localhost:8080"
    
    def test_health_check(self):
        """Test API gateway health endpoint"""
        response = requests.get(f"{self.BASE_URL}/healthz")
        assert response.status_code == 200
        assert response.json()["status"] == "healthy"
    
    def test_readiness_check(self):
        """Test API gateway readiness endpoint"""
        response = requests.get(f"{self.BASE_URL}/readyz")
        assert response.status_code == 200
        assert response.json()["ready"] == True
    
    def test_set_data(self):
        """Test setting data through API gateway"""
        response = requests.post(
            f"{self.BASE_URL}/api/v1/data/test_key",
            json={"value": "test_value"},
            headers={"Content-Type": "application/json"}
        )
        assert response.status_code == 200
        assert response.json()["success"] == True
    
    def test_get_data(self):
        """Test getting data through API gateway"""
        # First set data
        requests.post(
            f"{self.BASE_URL}/api/v1/data/test_get",
            json={"value": "get_value"},
            headers={"Content-Type": "application/json"}
        )
        
        # Then get data
        response = requests.get(f"{self.BASE_URL}/api/v1/data/test_get")
        assert response.status_code == 200
        assert response.json()["value"] == "get_value"
    
    def test_delete_data(self):
        """Test deleting data through API gateway"""
        # First set data
        requests.post(
            f"{self.BASE_URL}/api/v1/data/test_delete",
            json={"value": "delete_value"},
            headers={"Content-Type": "application/json"}
        )
        
        # Then delete data
        response = requests.delete(f"{self.BASE_URL}/api/v1/data/test_delete")
        assert response.status_code == 200
        assert response.json()["success"] == True


class TestRustDataStore:
    """Integration tests for Rust Data Store via gRPC"""
    
    GRPC_ADDRESS = "localhost:50051"
    
    @pytest.fixture
    def grpc_channel(self):
        """Create gRPC channel for testing"""
        channel = grpc.insecure_channel(self.GRPC_ADDRESS)
        yield channel
        channel.close()
    
    @pytest.fixture
    def grpc_client(self, grpc_channel):
        """Create gRPC client for testing"""
        client = data_store_pb2_grpc.DataStoreServiceClient(grpc_channel)
        yield client
    
    def test_health_check(self, grpc_client):
        """Test data store health check"""
        request = data_store_pb2.HealthCheckRequest()
        response = grpc_client.HealthCheck(request)
        assert response.healthy == True
    
    def test_set_operation(self, grpc_client):
        """Test set operation"""
        request = data_store_pb2.SetRequest(
            key=b"test_key",
            value=b"test_value",
            ttl_seconds=300
        )
        response = grpc_client.Set(request)
        assert response.success == True
    
    def test_get_operation(self, grpc_client):
        """Test get operation"""
        # First set data
        grpc_client.Set(data_store_pb2.SetRequest(
            key=b"test_get",
            value=b"get_value"
        ))
        
        # Then get data
        request = data_store_pb2.GetRequest(key=b"test_get")
        response = grpc_client.Get(request)
        assert response.found == True
        assert response.value == b"get_value"
    
    def test_delete_operation(self, grpc_client):
        """Test delete operation"""
        # First set data
        grpc_client.Set(data_store_pb2.SetRequest(
            key=b"test_delete",
            value=b"delete_value"
        ))
        
        # Then delete data
        request = data_store_pb2.DeleteRequest(key=b"test_delete")
        response = grpc_client.Delete(request)
        assert response.success == True


class TestPythonTelemetry:
    """Integration tests for Python Telemetry Collector"""
    
    BASE_URL = "http://localhost:8000"
    
    def test_health_check(self):
        """Test telemetry health endpoint"""
        response = requests.get(f"{self.BASE_URL}/healthz")
        assert response.status_code == 200
        assert response.json()["status"] == "healthy"
    
    def test_readiness_check(self):
        """Test telemetry readiness endpoint"""
        response = requests.get(f"{self.BASE_URL}/readyz")
        assert response.status_code == 200
        assert response.json()["ready"] == True
    
    def test_log_ingestion(self):
        """Test log ingestion endpoint"""
        log_data = {
            "service": "test-service",
            "level": "info",
            "message": "Test log message",
            "timestamp": "2024-01-01T00:00:00Z",
            "metadata": {
                "request_id": "test-123"
            }
        }
        response = requests.post(
            f"{self.BASE_URL}/api/v1/logs",
            json=log_data,
            headers={"Content-Type": "application/json"}
        )
        assert response.status_code == 200
        assert response.json()["success"] == True
    
    def test_metrics_query(self):
        """Test metrics query endpoint"""
        response = requests.get(
            f"{self.BASE_URL}/api/v1/metrics",
            params={"service": "api-gateway", "time_range": "1h"}
        )
        assert response.status_code == 200
        assert "metrics" in response.json()
    
    def test_anomaly_detection(self):
        """Test anomaly detection endpoint"""
        response = requests.get(f"{self.BASE_URL}/api/v1/anomalies")
        assert response.status_code == 200
        assert "anomalies" in response.json()


class TestEndToEnd:
    """End-to-end integration tests"""
    
    def test_complete_flow(self):
        """Test complete data flow through all services"""
        # 1. Set data via API Gateway
        gateway_response = requests.post(
            "http://localhost:8080/api/v1/data/e2e_test",
            json={"value": "e2e_value"},
            headers={"Content-Type": "application/json"}
        )
        assert gateway_response.status_code == 200
        
        # 2. Get data via API Gateway
        get_response = requests.get("http://localhost:8080/api/v1/data/e2e_test")
        assert get_response.status_code == 200
        assert get_response.json()["value"] == "e2e_value"
        
        # 3. Send log to telemetry
        log_response = requests.post(
            "http://localhost:8000/api/v1/logs",
            json={
                "service": "e2e-test",
                "level": "info",
                "message": "E2E test completed",
                "timestamp": "2024-01-01T00:00:00Z"
            },
            headers={"Content-Type": "application/json"}
        )
        assert log_response.status_code == 200
        
        # 4. Clean up
        delete_response = requests.delete("http://localhost:8080/api/v1/data/e2e_test")
        assert delete_response.status_code == 200


@pytest.fixture(scope="session", autouse=True)
def setup_services():
    """Setup services before running integration tests"""
    # This fixture can be used to start services if needed
    # For now, assume services are already running
    yield
    # Cleanup after tests
    pass


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
