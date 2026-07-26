#!/bin/bash

# Load Images Script
# This script loads Docker images into the Kind cluster

set -e

GREEN='\033[0;32m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

echo "Loading Docker images into Kind cluster..."

# Build images first
echo "Building Docker images..."
docker-compose build

# Load images into Kind
echo "Loading images into Kind cluster..."
kind load docker-image go-api-gateway:latest --name ci-cd-pipeline
print_success "Loaded go-api-gateway image"

kind load docker-image rust-data-store:latest --name ci-cd-pipeline
print_success "Loaded rust-data-store image"

kind load docker-image python-telemetry:latest --name ci-cd-pipeline
print_success "Loaded python-telemetry image"

print_success "All images loaded successfully!"
