#!/bin/bash

# MCP Server Installation Script for Kind Cluster
# This script installs the Kubernetes MCP Server with GitHub Copilot integration

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

echo "Installing MCP Server for GitHub Copilot integration..."

# Create namespace
kubectl create namespace mcp-server || print_warning "mcp-server namespace may already exist"

# Create service account and RBAC
echo "Setting up RBAC for MCP Server..."
kubectl apply -f manifests/bleeding-edge/mcp/rbac.yaml

print_success "RBAC configured"

# Apply MCP tool configurations
echo "Applying MCP tool configurations..."
kubectl apply -f manifests/bleeding-edge/mcp/config/

print_success "MCP tools configured"

# Deploy MCP Server
echo "Deploying MCP Server..."
kubectl apply -f manifests/bleeding-edge/mcp/mcp-server.yaml

print_success "MCP Server deployed"

# Apply GitHub Copilot integration
echo "Configuring GitHub Copilot integration..."
kubectl apply -f manifests/bleeding-edge/mcp/copilot-integration.yaml

print_success "GitHub Copilot integration configured"

# Wait for MCP Server to be ready
echo "Waiting for MCP Server to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/mcp-server -n mcp-server

print_success "MCP Server is ready!"

echo ""
echo "MCP Server is running with the following capabilities:"
echo "Read: pod logs, Hubble network flows, Prometheus metrics, deployment statuses"
echo "Write: apply manifests, trigger Argo CD syncs, scale deployments, spin up vClusters"
echo ""
echo "The MCP Server integrates with GitHub Copilot via webhook."
echo "Configure your GitHub Copilot settings to point to the MCP Server webhook endpoint."
