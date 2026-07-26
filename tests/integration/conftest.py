"""
Configuration and fixtures for integration tests
"""

import pytest
import subprocess
import time
import requests
import grpc
import os


@pytest.fixture(scope="session")
def service_urls():
    """Provide service URLs for testing"""
    return {
        "gateway": os.getenv("GATEWAY_URL", "http://localhost:8080"),
        "telemetry": os.getenv("TELEMETRY_URL", "http://localhost:8000"),
        "data_store": os.getenv("DATA_STORE_URL", "localhost:50051")
    }


@pytest.fixture(scope="session")
def wait_for_services(service_urls):
    """Wait for all services to be ready before running tests"""
    max_retries = 30
    retry_delay = 2
    
    print("Waiting for services to be ready...")
    
    # Wait for API Gateway
    for _ in range(max_retries):
        try:
            response = requests.get(f"{service_urls['gateway']}/healthz", timeout=5)
            if response.status_code == 200:
                print("API Gateway is ready")
                break
        except requests.exceptions.RequestException:
            time.sleep(retry_delay)
    else:
        pytest.fail("API Gateway did not become ready in time")
    
    # Wait for Telemetry
    for _ in range(max_retries):
        try:
            response = requests.get(f"{service_urls['telemetry']}/healthz", timeout=5)
            if response.status_code == 200:
                print("Telemetry is ready")
                break
        except requests.exceptions.RequestException:
            time.sleep(retry_delay)
    else:
        pytest.fail("Telemetry did not become ready in time")
    
    # Wait for Data Store
    for _ in range(max_retries):
        try:
            channel = grpc.insecure_channel(service_urls['data_store'])
            from proto import data_store_pb2, data_store_pb2_grpc
            client = data_store_pb2_grpc.DataStoreServiceClient(channel)
            response = client.HealthCheck(data_store_pb2.HealthCheckRequest(), timeout=5)
            if response.healthy:
                print("Data Store is ready")
                channel.close()
                break
        except Exception:
            time.sleep(retry_delay)
    else:
        pytest.fail("Data Store did not become ready in time")
    
    print("All services are ready!")
    yield
    
    print("Integration tests completed")


@pytest.fixture
def cleanup_data(service_urls):
    """Cleanup test data after each test"""
    yield
    # Clean up test data
    try:
        requests.delete(f"{service_urls['gateway']}/api/v1/data/test_*")
    except requests.exceptions.RequestException:
        pass
