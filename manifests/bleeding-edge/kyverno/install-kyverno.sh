#!/bin/bash

# Kyverno Installation Script for Kind Cluster
# This script installs Kyverno admission controller with security policies

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

echo "Installing Kyverno admission controller for Kind cluster..."

# Add Kyverno Helm repository
helm repo add kyverno https://kyverno.github.io/kyverno
helm repo update

# Create namespace
kubectl create namespace kyverno || print_warning "kyverno namespace may already exist"

# Install Kyverno in audit mode initially
helm install kyverno kyverno/kyverno \
    --namespace=kyverno \
    --version=3.1.0 \
    --set replicaCount=1 \
    --set admissionController.replicas=1 \
    --set backgroundController.replicas=1 \
    --set reportsController.replicas=1 \
    --set cleanupController.replicas=1

print_success "Kyverno installed successfully in audit mode!"

# Wait for Kyverno to be ready
echo "Waiting for Kyverno to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/kyverno-admission-controller -n kyverno
kubectl wait --for=condition=available --timeout=300s deployment/kyverno-background-controller -n kyverno

print_success "Kyverno is ready!"

# Apply security policies
echo "Applying Kyverno security policies..."
kubectl apply -f manifests/bleeding-edge/kyverno/policies/

print_success "Kyverno policies applied successfully!"

echo ""
echo "Kyverno is running in audit mode. Policies will report violations but not block deployments."
echo "To switch to enforce mode, update policy files and set validationFailureAction: Enforce"
echo ""
echo "View policy reports:"
echo "kubectl get cpol -A"
echo "kubectl get pol -A"
