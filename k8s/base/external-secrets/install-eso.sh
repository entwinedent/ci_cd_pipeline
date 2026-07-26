#!/bin/bash

# External Secrets Operator Installation Script
# This script installs ESO and configures it for Vault integration

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

echo "Installing External Secrets Operator..."

# Install ESO using Helm
helm repo add external-secrets https://charts.external-secrets.io
helm repo update

# Create namespace
kubectl create namespace external-secrets || print_warning "external-secrets namespace may already exist"

# Install ESO
helm install external-secrets external-secrets/external-secrets \
    --namespace=external-secrets \
    --version=0.9.0

print_success "External Secrets Operator installed successfully!"

# Wait for ESO to be ready
echo "Waiting for External Secrets Operator to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/external-secrets-operator -n external-secrets

print_success "External Secrets Operator is ready!"

# Apply custom manifests
echo "Applying ESO configuration..."
kubectl apply -f k8s/base/external-secrets-operator/
kubectl apply -f k8s/base/external-secrets/

print_success "ESO configured successfully!"
