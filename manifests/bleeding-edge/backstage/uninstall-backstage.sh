#!/bin/bash

# Backstage Uninstallation Script
# This script removes Backstage and all configurations

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

echo "Uninstalling Backstage..."

# Delete catalog and configuration
echo "Removing Backstage catalog and configuration..."
kubectl delete -f manifests/bleeding-edge/backstage/catalog-info.yaml --ignore-not-found=true
kubectl delete -f manifests/bleeding-edge/backstage/app-config.yaml --ignore-not-found=true

print_success "Backstage configuration removed"

# Uninstall Backstage
echo "Uninstalling Backstage Helm chart..."
helm uninstall backstage -n backstage || print_warning "Backstage may not be installed"

# Delete namespace
echo "Deleting backstage namespace..."
kubectl delete namespace backstage --ignore-not-found=true

print_success "Backstage uninstalled successfully!"
