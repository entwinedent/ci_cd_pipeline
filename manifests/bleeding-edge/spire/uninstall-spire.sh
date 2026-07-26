#!/bin/bash

# SPIRE Uninstallation Script
# This script removes SPIRE and all configurations

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

echo "Uninstalling SPIRE and configurations..."

# Delete workload entries
echo "Removing SPIRE workload entries..."
kubectl delete -f manifests/bleeding-edge/spire/entries/ --ignore-not-found=true

print_success "SPIRE workload entries removed"

# Delete bundle
echo "Removing SPIRE bundle..."
kubectl delete -f manifests/bleeding-edge/spire/bundle.yaml --ignore-not-found=true

print_success "SPIRE bundle removed"

# Delete SPIRE agent
echo "Removing SPIRE agent..."
kubectl delete -f manifests/bleeding-edge/spire/spire-agent.yaml --ignore-not-found=true

print_success "SPIRE agent removed"

# Delete SPIRE server
echo "Removing SPIRE server..."
kubectl delete -f manifests/bleeding-edge/spire/spire-server.yaml --ignore-not-found=true

print_success "SPIRE server removed"

# Uninstall SPIRE Helm chart
echo "Uninstalling SPIRE Helm chart..."
helm uninstall spire-server -n spire || print_warning "SPIRE Helm chart may not be installed"

print_success "SPIRE Helm chart removed"

# Delete namespace
echo "Deleting spire namespace..."
kubectl delete namespace spire --ignore-not-found=true

print_success "SPIRE uninstalled successfully!"
