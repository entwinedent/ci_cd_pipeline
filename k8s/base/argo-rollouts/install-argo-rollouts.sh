#!/bin/bash

# Argo Rollouts Installation Script for Kind Cluster
# This script installs Argo Rollouts for progressive delivery

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

echo "Installing Argo Rollouts for Kind cluster..."

# Create namespace
kubectl create namespace argo-rollouts || print_warning "argo-rollouts namespace may already exist"

# Install Argo Rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml

print_success "Argo Rollouts installed successfully!"

# Wait for Argo Rollouts to be ready
echo "Waiting for Argo Rollouts to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/argo-rollouts -n argo-rollouts

print_success "Argo Rollouts is ready!"

# Apply AnalysisTemplates
echo "Applying AnalysisTemplates..."
kubectl apply -f k8s/base/analysis-templates/

print_success "AnalysisTemplates applied successfully!"

# Display access info
echo ""
echo "Argo Rollouts Dashboard:"
echo "Port forward: kubectl port-forward svc/argo-rollouts-dashboard -n argo-rollouts 3100:3100"
echo "Access: http://localhost:3100"
echo ""
