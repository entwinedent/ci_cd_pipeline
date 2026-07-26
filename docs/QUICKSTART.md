# Quick Start Guide

This guide provides step-by-step instructions to get the CI/CD pipeline running locally in under 10 minutes.

## Prerequisites

- Docker Desktop (or Docker Engine)
- Kind (Kubernetes in Docker)
- kubectl
- Go 1.22+
- Rust (via rustup)
- Python

## Installation

### 1. Install Tools

```bash
# Install Kind
go install sigs.k8s.io/kind@latest

# Install kubectl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Install Python dependencies
pip install fastapi uvicorn pydantic httpx pytest pytest-cov
```

### 2. Clone Repository

```bash
git clone https://github.com/username/ci-cd-pipeline.git
cd ci-cd-pipeline
```

## Local Development

### Option 1: Docker Compose (Fastest)

```bash
# Start all services
docker-compose up

# Services available at:
# - Go API Gateway: http://localhost:8080
# - Python Telemetry: http://localhost:8000
# - Rust Data Store: localhost:50051
```

### Option 2: Kind Cluster (Full Kubernetes)

```bash
# Create Kind cluster
kind create cluster --config=config/kind-config.yaml

# Build Docker images
make docker-build

# Load images into cluster
make kind-setup

# Deploy services
make deploy

# Check status
make status

# Port forward to test
make port-forward
```

## Verification

### Test Services

```bash
# Test Go API Gateway
curl http://localhost:8080/healthz
curl -X POST http://localhost:8080/api/v1/data/test -d '{"value":"hello"}'
curl http://localhost:8080/api/v1/data/test

# Test Python Telemetry
curl http://localhost:8000/healthz
curl -X POST http://localhost:8000/api/v1/logs -H "Content-Type: application/json" -d '{"service":"test","level":"info","message":"test"}'
```

### Run Tests

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run smoke tests
make test-smoke
```

## Multi-Region DR Simulation

```bash
# Setup multi-region Kind clusters
bash scripts/dr-failover/setup-multi-kind.sh

# Simulate failover
bash scripts/dr-failover/simulate-failover.sh

# Check cluster health
bash scripts/dr-failover/simulate-failover.sh check

# Restore primary
bash scripts/dr-failover/simulate-failover.sh restore
```

## Cost Analysis

```bash
# Analyze Terraform costs
cd terraform
terraform plan -out=tfplan
terraform show -json tfplan > /tmp/plan.json
python ../tools/infracost-mock/infracost_mock.py --path /tmp/plan.json

# Analyze Kubernetes costs
python ../tools/kubectl-cost-mock/kubectl_cost_mock.py --path k8s/base
```

## Cleanup

```bash
# Stop Docker Compose
docker-compose down

# Delete Kind cluster
kind delete cluster --name ci-cd-pipeline

# Clean build artifacts
make clean
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :8080
lsof -i :8000

# Kill the process
kill -9 <PID>
```

### Docker Issues

```bash
# Restart Docker
sudo systemctl restart docker

# Check Docker status
docker ps
```

### Kind Cluster Issues

```bash
# Check cluster status
kind get clusters
kubectl cluster-info

# Delete and recreate
kind delete cluster --name ci-cd-pipeline
kind create cluster --config=config/kind-config.yaml
```

## Next Steps

- Read the [Architecture Documentation](docs/architecture.md)
- Explore [Service READMEs](services/)
- Try [Chaos Engineering Tests](tests/chaos/)
- Review [Multi-Region DR Setup](manifests/multi-region/)
