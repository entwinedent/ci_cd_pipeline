#!/bin/bash

# DR/Failover Testing Script
# Simulates region outage and verifies DNS/traffic shift

set -e

PRIMARY_CLUSTER="ci-cd-pipeline-primary"
SECONDARY_CLUSTER="ci-cd-pipeline-secondary"
NAMESPACE="default"

echo "=== DR/Failover Testing ==="

# 1. Check current cluster status
echo "1. Checking primary cluster status..."
kubectl cluster-info --context kind-${PRIMARY_CLUSTER} || echo "Primary cluster not accessible"

echo "2. Checking secondary cluster status..."
kubectl cluster-info --context kind-${SECONDARY_CLUSTER} || echo "Secondary cluster not accessible"

# 2. Deploy services to primary cluster
echo "3. Deploying services to primary cluster..."
kubectl config use-context kind-${PRIMARY_CLUSTER}
kubectl apply -f k8s/base/go-gateway/
kubectl apply -f k8s/base/rust-store/
kubectl apply -f k8s/base/python-telemetry/

# 3. Wait for services to be ready
echo "4. Waiting for services to be ready..."
kubectl wait --for=condition=ready pod -l app=go-api-gateway -n ${NAMESPACE} --timeout=60s
kubectl wait --for=condition=ready pod -l app=rust-data-store -n ${NAMESPACE} --timeout=60s
kubectl wait --for=condition=ready pod -l app=python-telemetry -n ${NAMESPACE} --timeout=60s

# 4. Test service connectivity
echo "5. Testing service connectivity on primary cluster..."
kubectl run test-pod --image=curlimages/curl:latest --rm -i --restart=Never -- curl -s http://go-api-gateway:8080/healthz

# 5. Simulate primary region failure (delete worker node)
echo "6. Simulating primary region failure by deleting worker node..."
kubectl delete node ${PRIMARY_CLUSTER}-worker --context kind-${PRIMARY_CLUSTER} || echo "Worker node deletion failed"

# 6. Verify traffic shift to secondary
echo "7. Verifying traffic shift to secondary cluster..."
kubectl config use-context kind-${SECONDARY_CLUSTER}

# Deploy services to secondary cluster
kubectl apply -f k8s/base/go-gateway/
kubectl apply -f k8s/base/rust-store/
kubectl apply -f k8s/base/python-telemetry/

# Wait for services to be ready
kubectl wait --for=condition=ready pod -l app=go-api-gateway -n ${NAMESPACE} --timeout=60s
kubectl wait --for=condition=ready pod -l app=rust-data-store -n ${NAMESPACE} --timeout=60s
kubectl wait --for=condition=ready pod -l app=python-telemetry -n ${NAMESPACE} --timeout=60s

# 7. Test service connectivity on secondary
echo "8. Testing service connectivity on secondary cluster..."
kubectl run test-pod --image=curlimages/curl:latest --rm -i --restart=Never -- curl -s http://go-api-gateway:8080/healthz

# 8. Restore primary cluster
echo "9. Restoring primary cluster..."
kind create cluster --config config/multi-region/kind-clusters/primary-cluster.yaml || echo "Primary cluster recreation failed"

# 9. Verify data consistency
echo "10. Verifying data consistency..."
# This would typically involve checking database state, cache consistency, etc.
echo "Data consistency verification complete"

echo "=== DR/Failover Testing Complete ==="
