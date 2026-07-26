#!/bin/bash

# Grafana Tempo Installation Script for Kind Cluster
# This script installs Tempo for distributed tracing with OTel integration

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

echo "Installing Grafana Tempo for distributed tracing..."

# Add Grafana Helm repository
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

# Create namespace
kubectl create namespace tempo || print_warning "tempo namespace may already exist"

# Install Tempo
helm install tempo grafana/tempo \
    --namespace=tempo \
    --version=1.8.0 \
    --values=manifests/bleeding-edge/opentelemetry/tempo-values.yaml

print_success "Tempo installed successfully!"

# Wait for Tempo to be ready
echo "Waiting for Tempo to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/tempo -n tempo

print_success "Tempo is ready!"

echo ""
echo "Tempo is running with the following configuration:"
echo "- Distributed tracing backend"
echo "- Integration with Prometheus via Grafana Exemplars"
echo "- OTel gRPC endpoint: tempo.tempo.svc.cluster.local:4317"
echo ""
echo "Access Tempo UI:"
echo "Port forward: kubectl port-forward svc/tempo -n tempo 3200:3200"
echo "Access: http://localhost:3200"
