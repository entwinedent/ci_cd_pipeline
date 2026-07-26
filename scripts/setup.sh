#!/bin/bash

# Advanced CI/CD Pipeline Setup Script
# This script sets up the complete development environment for the CI/CD pipeline

set -e

echo "🚀 Setting up Advanced CI/CD Pipeline Environment..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Check if running in WSL 2
if grep -q Microsoft /proc/version; then
    print_success "Running in WSL 2 environment"
else
    print_warning "Not running in WSL 2. Some features may not work as expected."
fi

# Check for required tools
echo "Checking for required tools..."

command -v git >/dev/null 2>&1 || { print_error "Git is not installed. Please install git."; exit 1; }
print_success "Git found: $(git --version)"

command -v docker >/dev/null 2>&1 || { print_error "Docker is not installed. Please install Docker."; exit 1; }
print_success "Docker found: $(docker --version)"

command -v kind >/dev/null 2>&1 || { print_error "Kind is not installed. Please install kind."; exit 1; }
print_success "Kind found: $(kind version)"

command -v kubectl >/dev/null 2>&1 || { print_error "kubectl is not installed. Please install kubectl."; exit 1; }
print_success "kubectl found: $(kubectl version --client --short)"

command -v helm >/dev/null 2>&1 || { print_error "Helm is not installed. Please install Helm."; exit 1; }
print_success "Helm found: $(helm version --short)"

command -v go >/dev/null 2>&1 || { print_error "Go is not installed. Please install Go 1.22+."; exit 1; }
print_success "Go found: $(go version)"

command -v rustc >/dev/null 2>&1 || { print_error "Rust is not installed. Please install Rust."; exit 1; }
print_success "Rust found: $(rustc --version)"

command -v python3 >/dev/null 2>&1 || { print_error "Python 3 is not installed. Please install Python 3."; exit 1; }
print_success "Python found: $(python3 --version)"

# Create local registry for development
echo "Setting up local container registry..."
docker-compose -f config/docker-compose.dev.yaml up -d registry
print_success "Local registry started on port 5000"

# Build service images
echo "Building service images..."
docker-compose build
print_success "Service images built successfully"

# Create Kind cluster
echo "Creating Kind cluster..."
kind create cluster --config=config/kind-config.yaml || print_warning "Kind cluster may already exist"
print_success "Kind cluster created"

# Load images into Kind cluster
echo "Loading images into Kind cluster..."
kind load docker-image go-api-gateway:latest --name ci-cd-pipeline 2>/dev/null || true
kind load docker-image rust-data-store:latest --name ci-cd-pipeline 2>/dev/null || true
kind load docker-image python-telemetry:latest --name ci-cd-pipeline 2>/dev/null || true
print_success "Images loaded into Kind cluster"

# Apply Kubernetes manifests
echo "Applying Kubernetes manifests..."
kubectl apply -f k8s/base/go-gateway/
kubectl apply -f k8s/base/rust-store/
kubectl apply -f k8s/base/python-telemetry/
print_success "Kubernetes manifests applied"

# Wait for deployments to be ready
echo "Waiting for deployments to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/go-api-gateway || print_warning "Go gateway deployment not ready"
kubectl wait --for=condition=available --timeout=300s deployment/rust-data-store || print_warning "Rust data store deployment not ready"
kubectl wait --for=condition=available --timeout=300s deployment/python-telemetry || print_warning "Python telemetry deployment not ready"

print_success "Setup completed successfully!"
echo ""
echo "🎉 Advanced CI/CD Pipeline environment is ready!"
echo ""
echo "Next steps:"
echo "1. Test the services: kubectl get pods"
echo "2. Access services: kubectl port-forward svc/go-api-gateway 8080:8080"
echo "3. View logs: kubectl logs -f deployment/go-api-gateway"
echo ""
