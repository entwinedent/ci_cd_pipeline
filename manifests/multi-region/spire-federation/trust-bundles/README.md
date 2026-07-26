# SPIRE Federation Trust Bundles

This directory contains the trust bundles for SPIRE server federation across regions.

## Architecture

The SPIRE federation enables cross-region identity validation:

- **Primary Region (us-east-1)**: Issues SVID certificates for workloads
- **Secondary Region (eu-central-1)**: Validates SVID certificates from primary region
- **Bidirectional Federation**: Both regions can issue and validate certificates

## Trust Bundle Management

### Automatic Synchronization

Trust bundles are automatically synchronized between regions every 30 seconds.

### Manual Rotation

To rotate trust bundles:

1. Generate new trust bundle on primary SPIRE server:
```bash
kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server bundle show -format pem > us-east-1-bundle.pem
```

2. Update ConfigMap:
``` Bash
kubectl create configmap spire-federation-bundles \
  --from-file=us-east-1-bundle.pem \
  --from-file=eu-central-1-bundle.pem \
  --dry-run=client -o yaml | kubectl apply -f -
```

3. Restart SPIRE agents in both regions:
```bash
kubectl rollout restart daemonset/spire-agent -n spire
```

## Federation Validation

To validate federation is working:

```bash
# Check trust bundle status
kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server bundle show

# Verify federation endpoints
kubectl get clusterfederation -n spire

# Check agent registration
kubectl get spireid -n spire
```

## Security Considerations

- Trust bundles contain sensitive cryptographic material
- Rotate trust bundles regularly (recommended: every 90 days)
- Monitor bundle synchronization for failures
- Implement alerting for bundle expiration
