# Vault Configuration

This directory contains HashiCorp Vault configuration for secrets management in the CI/CD pipeline.

## Purpose

Vault provides:
- Secure secrets storage
- Dynamic secrets generation
- Secrets rotation
- Access control and auditing

## Contents

### Installation

- **install-vault.sh** - Script to install Vault in Kubernetes cluster
- **values.yaml** - Helm values for Vault configuration

## Usage

### Installation

```bash
# Install Vault
bash config/vault/install-vault.sh

# Verify installation
kubectl get pods -n vault
kubectl get svc -n vault
```

### Access Vault UI

```bash
# Port forward Vault
kubectl port-forward svc/vault -n vault 8200:8200

# Access at http://localhost:8200
```

### Initialize Vault

```bash
# Initialize Vault
kubectl exec -n vault vault-0 -- vault operator init

# Unseal Vault
kubectl exec -n vault vault-0 -- vault operator unseal

# Login
kubectl exec -n vault vault-0 -- vault login
```

## Configuration

### Secrets Engines

Configure secrets engines for:
- **KV Secrets** - Static key-value storage
- **Database** - Dynamic database credentials
- **PKI** - Certificate authority
- **Transit** - Encryption as a service

### Authentication Methods

Configure authentication:
- **Kubernetes** - Service account authentication
- **Token** - Vault token authentication
- **GitHub** - GitHub OAuth
- **AppRole** - Application authentication

### Access Policies

Define policies for:
- Service-specific access
- Environment-specific access
- Role-based access control
- Audit logging

## Secrets Management

### Static Secrets

Store static secrets in KV engine:

```bash
# Write secret
vault kv put secret/api-gateway \
  api-key="your-api-key" \
  db-password="your-password"

# Read secret
vault kv get secret/api-gateway
```

### Dynamic Secrets

Generate dynamic credentials:

```bash
# Generate database credentials
vault read database/creds/api-gateway

# Generate PKI certificate
vault write pki/issue/api-gateway \
  common_name="api-gateway.local"
```

### Service Integration

Services integrate with Vault using:
- Kubernetes authentication
- Sidecar containers
- Environment variables
- SDK integration

## Best Practices

### Security

- Use least privilege access
- Rotate secrets regularly
- Enable audit logging
- Use dynamic secrets when possible

### High Availability

- Deploy Vault in HA mode
- Use Raft storage backend
- Configure standby nodes
- Monitor cluster health

### Backup and Recovery

- Regular snapshots
- Test recovery procedures
- Document backup process
- Store backups securely

## Troubleshooting

### Vault Not Starting

Check:
- Pod status and logs
- Resource limits
- Storage configuration
- Network policies

### Authentication Issues

Verify:
- Auth method configuration
- Service account permissions
- Token validity
- Policy configuration

### Secrets Not Accessible

Check:
- Secrets engine status
- Policy permissions
- Path configuration
- Token permissions

## Integration

### Kubernetes

Vault integrates with Kubernetes via:
- Kubernetes authentication method
- Mutating webhook for secrets injection
- Service account annotations
- Pod identity

### CI/CD

Vault integrates with CI/CD for:
- Secret injection during builds
- Dynamic credential generation
- Environment-specific secrets
- Audit trail

## Documentation

- [Vault Documentation](https://www.vaultproject.io/docs)
- [Vault Best Practices](https://www.vaultproject.io/docs/operations/production)
- [Kubernetes Integration](https://www.vaultproject.io/docs/platform/k8s)
