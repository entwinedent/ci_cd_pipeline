#!/bin/bash

# Kyverno Uninstallation Script
# This script removes Kyverno and all policies

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

echo "Uninstalling Kyverno..."

# Delete policies
echo "Removing Kyverno policies..."
kubectl delete -f manifests/bleeding-edge/kyverno/policies/ --ignore-not-found=true

print_success "Kyverno policies removed"

# Uninstall Kyverno
echo "Uninstalling Kyverno Helm chart..."
helm uninstall kyverno -n kyverno || print_warning "Kyverno may not be installed"

# Delete namespace
echo "Deleting kyverno namespace..."
kubectl delete namespace kyverno --ignore-not-found=true

print_success "Kyverno uninstalled successfully!"
