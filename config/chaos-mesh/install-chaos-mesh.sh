#!/bin/bash

# Chaos Mesh Installation Script for Kind Cluster
# This script installs Chaos Mesh for resilience testing

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

echo "Installing Chaos Mesh for Kind cluster..."

# Install Chaos Mesh using Helm
helm repo add chaos-mesh https://charts.chaos-mesh.org
helm repo update

# Create namespace
kubectl create namespace chaos-mesh || print_warning "chaos-mesh namespace may already exist"

# Install Chaos Mesh
helm install chaos-mesh chaos-mesh/chaos-mesh \
    --namespace=chaos-mesh \
    --version=2.6.0 \
    --set chaosDaemon.runtime=containerd \
    --set chaosDaemon.socketPath=/run/containerd/containerd.sock \
    --set dashboard.create=true

print_success "Chaos Mesh installed successfully!"

# Wait for Chaos Mesh to be ready
echo "Waiting for Chaos Mesh to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/chaos-controller-manager -n chaos-mesh
kubectl wait --for=condition=available --timeout=300s deployment/chaos-daemon -n chaos-mesh

print_success "Chaos Mesh is ready!"

# Display dashboard access info
echo ""
echo "Chaos Mesh Dashboard:"
echo "Port forward: kubectl port-forward svc/chaos-dashboard -n chaos-mesh 2333:2333"
echo "Access: http://localhost:2333"
echo ""

# Apply chaos experiments
echo "Applying chaos experiment configurations..."
kubectl apply -f config/chaos-mesh/network-delay.yaml
kubectl apply -f config/chaos-mesh/network-delay-rust-store.yaml
kubectl apply -f config/chaos-mesh/pod-kill-api-gateway.yaml
kubectl apply -f config/chaos-mesh/pod-kill-rust-store.yaml
kubectl apply -f config/chaos-mesh/memory-stress-api-gateway.yaml
kubectl apply -f config/chaos-mesh/memory-stress-rust-store.yaml

print_success "Chaos experiments configured!"
echo ""
echo "Note: Chaos experiments are paused by default. Use Chaos Mesh Dashboard or kubectl to enable them."
