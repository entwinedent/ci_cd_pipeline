#!/bin/bash

# Cilium Installation Script for Kind Cluster
# This script installs Cilium with eBPF networking and Hubble observability

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

echo "Installing Cilium with eBPF networking for Kind cluster..."

# Add Cilium Helm repository
helm repo add cilium https://helm.cilium.io/
helm repo update

# Create namespace
kubectl create namespace cilium || print_warning "cilium namespace may already exist"

# Install Cilium in Kind-compatible mode
helm install cilium cilium/cilium \
    --namespace=cilium \
    --version=1.14.0 \
    --values=manifests/bleeding-edge/cilium/values.yaml

print_success "Cilium installed successfully!"

# Wait for Cilium to be ready
echo "Waiting for Cilium to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/cilium-operator -n cilium
kubectl wait --for=condition=available --timeout=300s deployment/cilium -n cilium

print_success "Cilium is ready!"

# Enable Hubble
echo "Enabling Hubble observability..."
bash manifests/bleeding-edge/cilium/enable-hubble.sh

print_success "Hubble enabled successfully!"

echo ""
echo "Cilium is running in complementary mode alongside Kind's default CNI."
echo "eBPF host routing is enabled for zero-overhead network tracing."
echo ""
echo "Access Hubble UI:"
echo "Port forward: kubectl port-forward svc/hubble-ui -n cilium 12000:80"
echo "Access: http://localhost:12000"
