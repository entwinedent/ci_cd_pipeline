#!/bin/bash

# Tempo Uninstallation Script
# This script removes Tempo and OTel configurations

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

echo "Uninstalling Tempo and OTel configurations..."

# Delete OTel configurations
echo "Removing OTel instrumentation configurations..."
kubectl delete -f manifests/bleeding-edge/opentelemetry/instrumentation/ --ignore-not-found=true

print_success "OTel configurations removed"

# Delete OTel collector
echo "Removing OTel collector..."
kubectl delete -f manifests/bleeding-edge/opentelemetry/otel-collector-sidecar.yaml --ignore-not-found=true

print_success "OTel collector removed"

# Uninstall Tempo
echo "Uninstalling Tempo Helm chart..."
helm uninstall tempo -n tempo || print_warning "Tempo may not be installed"

# Delete namespace
echo "Deleting tempo namespace..."
kubectl delete namespace tempo --ignore-not-found=true

print_success "Tempo uninstalled successfully!"
