#!/bin/bash

# SPIRE Installation Script for Kind Cluster
# This script installs SPIRE for workload identity and universal mTLS

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

echo "Installing SPIRE for workload identity and universal mTLS..."

# Add SPIRE Helm repository
helm repo add spire https://charts.spiffe.io
helm repo update

# Create namespace
kubectl create namespace spire || print_warning "spire namespace may already exist"

# Install SPIRE Server
helm install spire-server spire/spire-server \
    --namespace=spire \
    --version=0.20.0 \
    --values=manifests/bleeding-edge/spire/spire-server.yaml

print_success "SPIRE Server installed successfully!"

# Wait for SPIRE Server to be ready
echo "Waiting for SPIRE Server to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/spire-server -n spire

print_success "SPIRE Server is ready!"

# Install SPIRE Agent
kubectl apply -f manifests/bleeding-edge/spire/spire-agent.yaml

print_success "SPIRE Agent installed successfully!"

# Wait for SPIRE Agent to be ready
echo "Waiting for SPIRE Agent to be ready..."
kubectl wait --for=condition=ready --timeout=300s daemonset/spire-agent -n spire

print_success "SPIRE Agent is ready!"

# Apply trust bundle
kubectl apply -f manifests/bleeding-edge/spire/bundle.yaml

print_success "Trust bundle configured!"

# Register workload entries
kubectl apply -f manifests/bleeding-edge/spire/entries/

print_success "Workload entries registered!"

echo ""
echo "SPIRE is running with the following configuration:"
echo "- Universal mTLS for all service-to-service communication"
echo "- Short-lived X.509 certificate rotation"
echo "- Complements Cilium network policies"
echo "- Cryptographic workload identity (SVID certificates)"
echo ""
echo "SPIRE Server endpoint: spire-server.spire.svc.cluster.local:8081"
