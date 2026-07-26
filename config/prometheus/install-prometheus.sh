#!/bin/bash

# Prometheus Installation Script for Kind Cluster
# This script installs Prometheus for metrics collection

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

echo "Installing Prometheus for Kind cluster..."

# Add Prometheus Helm repository
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Create namespace
kubectl create namespace monitoring || print_warning "monitoring namespace may already exist"

# Install Prometheus
helm install prometheus prometheus-community/kube-prometheus-stack \
    --namespace=monitoring \
    --values=config/prometheus/values.yaml \
    --version=55.0.0 \
    --set grafana.enabled=true \
    --set grafana.adminPassword=admin

print_success "Prometheus installed successfully!"

# Wait for Prometheus to be ready
echo "Waiting for Prometheus to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/prometheus-kube-prometheus-prometheus -n monitoring
kubectl wait --for=condition=available --timeout=300s deployment/grafana -n monitoring

print_success "Prometheus and Grafana are ready!"

# Display access info
echo ""
echo "Prometheus UI:"
echo "Port forward: kubectl port-forward svc/prometheus-kube-prometheus-prometheus -n monitoring 9090:9090"
echo "Access: http://localhost:9090"
echo ""
echo "Grafana UI:"
echo "Port forward: kubectl port-forward svc/grafana -n monitoring 3000:3000"
echo "Access: http://localhost:3000"
echo "Username: admin"
echo "Password: admin"
echo ""
