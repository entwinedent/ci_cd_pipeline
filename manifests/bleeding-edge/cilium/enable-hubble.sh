#!/bin/bash

# Hubble Enablement Script
# This script enables Hubble CLI and configures flow collection

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

echo "Enabling Hubble observability..."

# Enable Hubble CLI
cilium hubble enable --ui

print_success "Hubble CLI enabled"

# Apply Hubble UI deployment
kubectl apply -f manifests/bleeding-edge/cilium/hubble-ui.yaml

print_success "Hubble UI deployed"

# Wait for Hubble components
echo "Waiting for Hubble components to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/hubble-ui -n cilium
kubectl wait --for=condition=available --timeout=300s deployment/hubble-relay -n cilium

print_success "Hubble is ready!"

echo ""
echo "Hubble CLI commands:"
echo "  cilium hubble flow show    # Show network flows"
echo "  cilium hubble status       # Check Hubble status"
echo "  cilium hubble port-forward # Port forward UI"
