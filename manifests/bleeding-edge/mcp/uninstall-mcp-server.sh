#!/bin/bash

# MCP Server Uninstallation Script
# This script removes the MCP Server and all configurations

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

echo "Uninstalling MCP Server..."

# Delete configurations
echo "Removing MCP tool configurations..."
kubectl delete -f manifests/bleeding-edge/mcp/config/ --ignore-not-found=true

print_success "MCP tool configurations removed"

# Delete GitHub Copilot integration
echo "Removing GitHub Copilot integration..."
kubectl delete -f manifests/bleeding-edge/mcp/copilot-integration.yaml --ignore-not-found=true

print_success "GitHub Copilot integration removed"

# Delete MCP Server deployment
echo "Removing MCP Server deployment..."
kubectl delete -f manifests/bleeding-edge/mcp/mcp-server.yaml --ignore-not-found=true

print_success "MCP Server removed"

# Delete RBAC
echo "Removing RBAC..."
kubectl delete -f manifests/bleeding-edge/mcp/rbac.yaml --ignore-not-found=true

print_success "RBAC removed"

# Delete namespace
echo "Deleting mcp-server namespace..."
kubectl delete namespace mcp-server --ignore-not-found=true

print_success "MCP Server uninstalled successfully!"
