# SPIFFE/SPIRE Zero-Trust Identity

SPIRE server and agent deployment for workload identity with X.509 SVID certificates and multi-region federation.

## Architecture

### Components

**SPIRE Server**
- Central certificate authority for the cluster
- Issues X.509 SVID certificates to workloads
- Trust bundle management and distribution
- Federation with other clusters

**SPIRE Agent**
- Runs as DaemonSet on each node
- Issues SVIDs to workloads on the node
- Rotates certificates automatically
- Enforces workload identity policies

### Trust Domain

- **Trust Domain**: `ci-cd-pipeline.local`
- **Bundle Endpoint**: SPIRE server API
- **Socket**: Unix domain socket for agent communication

## Deployment

### Deploy SPIRE

```bash
kubectl apply -f manifests/bleeding-edge/spire/
```

### Verify Deployment

```bash
# Check SPIRE server
kubectl get pods -n spire -l app=spire-server

# Check SPIRE agents
kubectl get pods -n spire -l app=spire-agent

# Check trust bundle
kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server bundle show
```

## Workload Identity

### Workload Entries

Each microservice has a SPIRE workload entry:

```yaml
apiVersion: spire.spiffe.io/v1alpha1
kind: ClusterWorkloadEntry
metadata:
  name: go-api-gateway
spec:
  selector:
    app: go-api-gateway
  spiffeID:
    trustDomain: ci-cd-pipeline.local
    path: /ns/default/sa/go-api-gateway
```

### SVID Certificate

Workloads receive X.509 SVID certificates:
- SPIFFE ID: `spiffe://ci-cd-pipeline.local/ns/default/sa/go-api-gateway`
- Valid for 1 hour (auto-rotated)
- Contains workload identity and trust domain

## Multi-Region Federation

### Federation Configuration

SPIRE servers federate across regions:

```yaml
apiVersion: spire.spiffe.io/v1alpha1
kind: ClusterFederation
metadata:
  name: multi-region-federation
spec:
  trustDomain: ci-cd-pipeline.local
  mode: bidirectional
```

### Trust Bundle Sync

Trust bundles are synchronized every 30 seconds:
- Primary region bundle shared with secondary
- Secondary region bundle shared with primary
- Automatic failover on bundle unavailability

## Configuration

### Environment Variables

- `SPIRE_SERVER_SOCKET` - SPIRE server socket path
- `SPIRE_AGENT_SOCKET` - SPIRE agent socket path
- `TRUST_DOMAIN` - SPIFFE trust domain

### Service Account

Workloads need service account annotation:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: go-api-gateway
  annotations:
    spiffe.io/spire-server: spire-server.spire.svc.cluster.local
```

## Security

- Zero-trust model with mutual TLS
- Automatic certificate rotation
- Fine-grained workload identity
- No static credentials or API keys
- Multi-region federation for cross-cluster trust

## Troubleshooting

### Workload Not Getting SVID

```bash
# Check agent logs
kubectl logs -n spire -l app=spire-agent

# Check workload entry
kubectl get clusterworkloadentry -n spire

# Verify service account annotation
kubectl get sa go-api-gateway -o yaml
```

### Federation Not Working

```bash
# Check federation status
kubectl get clusterfederation -n spire

# Check trust bundle sync
kubectl get configmap spire-federation-bundles -n spire

# Verify bundle endpoints
kubectl exec -n spire spire-server-0 -- /opt/spire/bin/spire-server bundle show
```

## Documentation

- [SPIRE Documentation](https://spiffe.io/docs/latest/spire/)
- [SPIFFE Specification](https://github.com/spiffe/spiffe/blob/main/standards/SPIFFE.md)
- [Multi-Cluster Federation](https://spiffe.io/docs/latest/spire/running-spire/federation/)
