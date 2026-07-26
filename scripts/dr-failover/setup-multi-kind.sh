#!/bin/bash

# Multi-Region Kind Cluster Setup Script
# Creates 3 Kind clusters to simulate multi-region DR architecture

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo "kind is not installed. Please install it first."
    exit 1
fi

# Create primary cluster (us-east-1 simulation)
print_info "Creating primary cluster (us-east-1 simulation)..."
kind create cluster --config=config/multi-region/kind-clusters/primary-cluster.yaml
print_success "Primary cluster created"

# Create secondary cluster (eu-central-1 simulation)
print_info "Creating secondary cluster (eu-central-1 simulation)..."
kind create cluster --config=config/multi-region/kind-clusters/secondary-cluster.yaml
print_success "Secondary cluster created"

# Create local dev cluster
print_info "Creating local development cluster..."
kind create cluster --config=config/multi-region/kind-clusters/local-dev-cluster.yaml
print_success "Local dev cluster created"

# Install Cilium on all clusters
print_info "Installing Cilium on primary cluster..."
kubectl config use-context kind-ci-cd-pipeline-primary
cilium install --version 1.14.0
print_success "Cilium installed on primary cluster"

print_info "Installing Cilium on secondary cluster..."
kubectl config use-context kind-ci-cd-pipeline-secondary
cilium install --version 1.14.0
print_success "Cilium installed on secondary cluster"

print_info "Installing Cilium on local dev cluster..."
kubectl config use-context kind-ci-cd-pipeline-local
cilium install --version 1.14.0
print_success "Cilium installed on local dev cluster"

# Set up cluster networking
print_info "Setting up multi-cluster networking..."
kubectl config use-context kind-ci-cd-pipeline-primary

# Create cluster context aliases
kubectl config set-context kind-ci-cd-pipeline-primary --cluster=kind-ci-cd-pipeline-primary
kubectl config set-context kind-ci-cd-pipeline-secondary --cluster=kind-ci-cd-pipeline-secondary
kubectl config set-context kind-ci-cd-pipeline-local --cluster=kind-ci-cd-pipeline-local

print_success "Multi-region Kind cluster setup complete!"
print_info "Cluster contexts:"
print_info "  - kind-ci-cd-pipeline-primary (us-east-1 simulation)"
print_info "  - kind-ci-cd-pipeline-secondary (eu-central-1 simulation)"
print_info "  - kind-ci-cd-pipeline-local (local development)"
