# Argo CD Multi-Cluster GitOps

Argo CD ApplicationSets for multi-region synchronization across us-east-1 and eu-central-1 clusters.

## Architecture

### Hub-and-Spoke Pattern

**Hub (Control Plane)**
- Centralized Argo CD deployment
- Manages multiple remote clusters
- Coordinates deployments across regions

**Spokes (Worker Clusters)**
- EKS clusters in us-east-1 and eu-central-1
- Receive deployment instructions from hub
- Maintain local cluster state

### ApplicationSets

ApplicationSets enable multi-cluster deployments:
- List generator for cluster list
- Template for deployment manifests
- Automated sync policies
- Self-healing and pruning

## Components

### Cluster Secrets

Each cluster has a secret for Argo CD access:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: region-us-east-1
  namespace: argocd
  labels:
    argocd.argoproj.io/secret-type: cluster
```

### ApplicationSets

**Multi-Region Workloads**
- Deploys core services to both regions
- Automated sync with self-healing
- Region-specific overlays

**Multi-Region Networking**
- Deploys Cilium to both regions
- Consistent network policies
- eBPF observability

**Multi-Region Security**
- Deploys SPIRE to both regions
- Identity federation
- Zero-trust networking

## Deployment

### Add Clusters to Argo CD

```bash
# Add us-east-1 cluster
argocd cluster add region-us-east-1

# Add eu-central-1 cluster
argocd cluster add region-eu-central-1

# Verify clusters
argocd cluster list
```

### Apply ApplicationSets

```bash
kubectl apply -f manifests/multi-region/argocd/applicationsets/
```

### Verify Deployment

```bash
# Check ApplicationSets
argocd appset list

# Check applications
argocd app list

# Check sync status
argocd app get core-service-us-east-1
argocd app get core-service-eu-central-1
```

## Configuration

### Cluster List

Edit the ApplicationSet to add/remove clusters:

```yaml
generators:
  - list:
      elements:
        - cluster: region-us-east-1
          url: https://gr7.us-east-1.eks.amazonaws.com
          region: us-east-1
        - cluster: region-eu-central-1
          url: https://gr7.eu-central-1.eks.amazonaws.com
          region: eu-central-1
```

### Sync Policy

Automated sync with self-healing:

```yaml
syncPolicy:
  automated:
    prune: true
    selfHeal: true
    allowEmpty: false
```

## Multi-Region Sync

### Sync Process

1. Developer pushes to Git
2. Argo CD detects change
3. ApplicationSet generates applications
4. Both clusters sync simultaneously
5. Argo CD monitors cluster health
6. Auto-remediation on failure

### Failover Handling

If primary region fails:
- Secondary region continues serving traffic
- Argo CD maintains desired state
- Manual intervention for failback
- GitOps ensures consistency

## Observability

### Argo CD Metrics

Prometheus metrics for:
- Application sync status
- Cluster health
- Deployment duration
- Error rates

### Alerts

Configure alerts for:
- Sync failures
- Cluster unavailability
- Drift detection
- Health check失败

## Troubleshooting

### Sync Failures

```bash
# Check application sync status
argocd app get <app-name>

# Check application logs
argocd app logs <app-name>

# Force sync
argocd app sync <app-name> --force
```

### Cluster Unavailable

```bash
# Check cluster status
argocd cluster get <cluster-name>

# Reconnect cluster
argocd cluster add <cluster-name>

# Check cluster credentials
kubectl get secret -n argocd <cluster-name>
```

## Documentation

- [Argo CD Documentation](https://argoproj.github.io/argo-cd/)
- [ApplicationSets Guide](https://argoproj.github.io/argo-cd/user-guide/applicationset/)
- [Multi-Cluster Management](https://argoproj.github.io/argo-cd/operator-manual/cluster-management/)
