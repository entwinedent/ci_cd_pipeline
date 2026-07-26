#!/bin/bash

# Cilium Uninstallation Script
# This script removes Cilium and Hubble

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

echo "Uninstalling Cilium..."

# Delete network policies
echo "Removing Cilium network policies..."
kubectl delete -f manifests/bleeding-edge/cilium/network-policies/ --ignore-not-found=true

print_success "Cilium network policies removed"

# Disable Hubble
echo "Disabling Hubble..."
cilium hubble disable || print_warning "Hubble may not be enabled"

print_success "Hubble disabled"

# Uninstall Cilium
echo "Uninstalling Cilium Helm chart..."
helm uninstall cilium -n cilium || print_warning "Cilium may not be installed"

# Delete namespace
echo "Deleting cilium namespace..."
kubectl delete namespace cilium --ignore-not-found=true

print_success "Cilium uninstalled successfully!"
