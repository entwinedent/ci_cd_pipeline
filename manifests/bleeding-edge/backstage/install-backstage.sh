#!/bin/bash

# Backstage Installation Script for Kind Cluster
# This script installs Backstage developer portal with Argo CD and vCluster integration

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

echo "Installing Backstage developer portal for Kind cluster..."

# Add Backstage Helm repository
helm repo add backstage https://backstage.github.io/backstage-chart
helm repo update

# Create namespace
kubectl create namespace backstage || print_warning "backstage namespace may already exist"

# Install Backstage
helm install backstage backstage/backstage \
    --namespace=backstage \
    --version=1.10.0 \
    --values=manifests/bleeding-edge/backstage/values.yaml

print_success "Backstage installed successfully!"

# Wait for Backstage to be ready
echo "Waiting for Backstage to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/backstage -n backstage

print_success "Backstage is ready!"

# Apply catalog and configuration
echo "Applying Backstage catalog and configuration..."
kubectl apply -f manifests/bleeding-edge/backstage/catalog-info.yaml
kubectl apply -f manifests/bleeding-edge/backstage/app-config.yaml

print_success "Backstage configuration applied!"

echo ""
echo "Backstage is running with the following integrations:"
echo "- Argo CD plugin for deployment visibility"
echo "- Kubernetes plugin for cluster resources"
echo "- vCluster plugin for ephemeral environments"
echo "- GitHub plugin for repository catalog"
echo ""
echo "Access Backstage:"
echo "Port forward: kubectl port-forward svc/backstage -n backstage 3000:80"
echo "Access: http://localhost:3000"
