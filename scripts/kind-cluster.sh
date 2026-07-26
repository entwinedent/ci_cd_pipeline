#!/bin/bash

# Kind Cluster Management Script
# This script manages the local Kind cluster for development

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

case "${1:-}" in
  create)
    echo "Creating Kind cluster..."
    kind create cluster --config=config/kind-config.yaml
    print_success "Kind cluster created"
    ;;
  
  delete)
    echo "Deleting Kind cluster..."
    kind delete cluster --name ci-cd-pipeline
    print_success "Kind cluster deleted"
    ;;
  
  status)
    echo "Kind cluster status:"
    kind get clusters
    kubectl cluster-info --context kind-ci-cd-pipeline
    ;;
  
  load-images)
    echo "Loading images into Kind cluster..."
    kind load docker-image go-api-gateway:latest --name ci-cd-pipeline
    kind load docker-image rust-data-store:latest --name ci-cd-pipeline
    kind load docker-image python-telemetry:latest --name ci-cd-pipeline
    print_success "Images loaded into Kind cluster"
    ;;
  
  apply-manifests)
    echo "Applying Kubernetes manifests..."
    kubectl apply -f k8s/base/go-gateway/
    kubectl apply -f k8s/base/rust-store/
    kubectl apply -f k8s/base/python-telemetry/
    print_success "Kubernetes manifests applied"
    ;;
  
  *)
    echo "Usage: $0 {create|delete|status|load-images|apply-manifests}"
    echo ""
    echo "Commands:"
    echo "  create           - Create Kind cluster"
    echo "  delete           - Delete Kind cluster"
    echo "  status           - Show cluster status"
    echo "  load-images      - Load Docker images into cluster"
    echo "  apply-manifests   - Apply Kubernetes manifests"
    exit 1
    ;;
esac
