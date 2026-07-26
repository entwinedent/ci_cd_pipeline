# Platform Deployment Guide

## Prerequisites

### Kind Installation
Kind is required for local Kubernetes cluster deployment. Install it from: https://kind.sigs.k8s.io/

**Windows Installation:**
```powershell
choco install kind
```

**Linux/WSL Installation:**
```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind
```

### Additional Tools
- Helm 3.x
- kubectl
- Docker or Docker Desktop

## Deployment Scripts

The following scripts are available in `manifests/bleeding-edge/`:

### 1. Argo CD
```bash
bash scripts/install-argocd.sh
```

### 2. Cilium + Hubble
```bash
bash manifests/bleeding-edge/cilium/install-cilium.sh
```

### 3. SPIRE
```bash
bash manifests/bleeding-edge/spire/install-spire.sh
```

### 4. Backstage
```bash
bash manifests/bleeding-edge/backstage/install-backstage.sh
```

### 5. Unleash
```bash
kubectl apply -f manifests/bleeding-edge/unleash/unleash-deployment.yaml
kubectl apply -f manifests/bleeding-edge/unleash/unleash-config.yaml
```

## Deployment Order

1. Create Kind cluster: `kind create cluster --config config/kind-config.yaml`
2. Install Cilium + Hubble (networking + observability)
3. Install SPIRE (zero-trust identity)
4. Install Argo CD (GitOps)
5. Install Backstage (developer portal)
6. Install Unleash (feature flags)
7. Deploy microservices via Argo CD

## Current Status

- ✅ Docker Compose services running (go-api-gateway, python-telemetry, rust-data-store)
- ❌ Kind not installed on this Windows system
- ✅ Installation scripts available for all platform services
- ✅ Screenshot documentation created

## Alternative Approach

Since Kind is not available, screenshots can be captured from:
1. GitHub Actions workflows (no Kind required)
2. Docker Compose service logs and health checks
3. Mock screenshots using the existing installation scripts as reference
