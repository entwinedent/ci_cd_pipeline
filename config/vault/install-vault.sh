#!/bin/bash

# HashiCorp Vault Installation Script for Kind Cluster
# This script installs Vault and configures it for the CI/CD pipeline

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

echo "Installing HashiCorp Vault for Kind cluster..."

# Add HashiCorp Helm repository
helm repo add hashicorp https://helm.releases.hashicorp.com
helm repo update

# Create namespace
kubectl create namespace vault || print_warning "vault namespace may already exist"

# Install Vault using Helm
helm install vault hashicorp/vault \
    --namespace=vault \
    --values=config/vault/values.yaml \
    --version=0.25.0

print_success "Vault installed successfully!"

# Wait for Vault to be ready
echo "Waiting for Vault to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/vault -n vault

print_success "Vault is ready!"

# Initialize Vault (dev mode is already initialized)
echo "Configuring Vault..."

# Enable KV secrets engine
kubectl exec -n vault vault-0 -- vault secrets enable -path=ci-cd-pipeline kv-v2

# Create sample secrets
kubectl exec -n vault vault-0 -- vault kv put ci-cd-pipeline/api-gateway \
    log_level=info \
    data_store_target=rust-data-store:50051

kubectl exec -n vault vault-0 -- vault kv put ci-cd-pipeline/python-telemetry \
    log_level=info \
    argocd_webhook_url="" \
    slack_webhook_url=""

# Create policy for External Secrets Operator
kubectl exec -n vault vault-0 -- vault policy write ci-cd-pipeline-policy - <<EOF
path "ci-cd-pipeline/data/*" {
  capabilities = ["read", "list"]
}
EOF

print_success "Vault configured successfully!"

# Display Vault access info
echo ""
echo "Vault UI:"
echo "Port forward: kubectl port-forward svc/vault -n vault 8200:8200"
echo "Access: http://localhost:8200"
echo "Token: Use root token from Vault pod (dev mode)"
echo ""
echo "Get root token:"
echo "kubectl exec -n vault vault-0 -- vault print token"
