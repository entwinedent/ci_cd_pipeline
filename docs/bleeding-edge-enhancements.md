# Bleeding-Edge CI/CD Pipeline Enhancements

This document describes the cutting-edge platform engineering components added to elevate the CI/CD pipeline to top-tier engineering standards, including Kyverno admission control, Cilium eBPF networking, Backstage developer portal, and MCP AI-Ops integration.

## Overview

These four bleeding-edge additions represent the frontier of modern platform engineering, used by top-tier organizations to achieve:

- **Zero-trust security** at the Kubernetes admission layer
- **Kernel-level observability** without sidecar overhead
- **Developer self-service** through a unified portal
- **AI-powered operations** via GitHub Copilot integration

## Phase 1: Kyverno Admission Control

### Location
`manifests/bleeding-edge/kyverno/`

### Features
- **Native Kubernetes YAML-based policies** (no Rego complexity)
- **Audit mode initially** for safe rollout
- **Comprehensive security policies**:
  - Rootless container enforcement
  - Required labels (app, version, environment)
  - Cosign image signature verification
  - Resource requests/limits enforcement
  - Network policy requirements

### Policies Implemented

1. **rootless-containers.yaml**
   - Enforces non-zero user ID for all containers
   - Requires read-only root filesystem
   - Prevents privilege escalation attacks

2. **required-labels.yaml**
   - Requires standard labels on all workloads
   - Ensures proper organization and tracking
   - Applies to Pods, Deployments, StatefulSets, DaemonSets

3. **cosign-verification.yaml**
   - Verifies Cosign image signatures before deployment
   - Uses Sigstore keyless verification
   - Ensures only signed images from trusted sources run

4. **resource-limits.yaml**
   - Enforces resource requests (CPU, memory)
   - Enforces resource limits (CPU, memory)
   - Prevents resource starvation

#### Installation
```bash
bash manifests/bleeding-edge/kyverno/install-kyverno.sh
```

#### Uninstallation
```bash
bash manifests/bleeding-edge/kyverno/uninstall-kyverno.sh
```

#### Usage
```bash
# View policy reports
kubectl get cpol -A
kubectl get pol -A

# Switch to enforce mode
# Edit policy files and set validationFailureAction: Enforce
```

## Phase 2: Cilium + Hubble eBPF Networking

### Location
`manifests/bleeding-edge/cilium/`

### Features
- **Complementary CNI mode** (alongside Kind's default CNI)
- **eBPF host routing** for zero-overhead network tracing
- **Hubble observability** without sidecar injection
- **L7 network policies** for service-to-service communication
- **Real-time flow visualization** via Hubble UI

### Components

1. **Cilium Installation**
   - eBPF-based networking
   - Host routing enabled
   - Kube-proxy replacement (partial mode)

2. **Hubble Observability**
   - Real-time network flow logs
   - Distributed tracing
   - Service-to-service communication visualization
   - Metrics integration with Prometheus

3. **L7 Network Policies**
   - API Gateway → Data Store (gRPC)
   - API Gateway → Telemetry (HTTP)
   - DNS access for all services

### Installation
```bash
bash manifests/bleeding-edge/cilium/install-cilium.sh
```

### Uninstallation
```bash
bash manifests/bleeding-edge/cilium/uninstall-cilium.sh
```

### Usage
```bash
# Access Hubble UI
kubectl port-forward svc/hubble-ui -n cilium 12000:80
# Open: http://localhost:12000

# View network flows
cilium hubble flow show

# Check Hubble status
cilium hubble status
```

## Phase 3: Backstage Developer Portal

### Location
`manifests/bleeding-edge/backstage/`

### Features
- **CNCF graduated** developer portal
- **Software Templates (Scaffolder)** for Go, Rust, Python
- **Argo CD plugin** for deployment visibility
- **Kubernetes plugin** for cluster resources
- **vCluster plugin** for ephemeral environments
- **GitHub integration** for repository catalog

### Components

1. **Backstage Installation**
   - PostgreSQL database
   - Plugin configuration
   - Service account for Kubernetes access

2. **Software Templates**
   - Go microservice template
   - Rust microservice template
   - Python microservice template
   - Each includes: Dockerfile, Helm chart, CI/CD pipeline

3. **Catalog Configuration**
   - Component catalog (API Gateway, Data Store, Telemetry)
   - System catalog (CI/CD Pipeline)
   - Resource catalog (Argo CD)

### Installation
```bash
bash manifests/bleeding-edge/backstage/install-backstage.sh
```

### Uninstallation
```bash
bash manifests/bleeding-edge/backstage/uninstall-backstage.sh
```

### Usage
```bash
# Access Backstage
kubectl port-forward svc/backstage -n backstage 3000:80
# Open: http://localhost:3000

# Configure environment variables
export GITHUB_TOKEN=your_github_token
export ARGOCD_USERNAME=admin
export ARGOCD_PASSWORD=your_argocd_password
```

### Creating a New Service
1. Navigate to "Create" in Backstage
2. Select "Go Microservice Template" (or Rust/Python)
3. Fill in service details
4. Backstage scaffolds the complete service
5. Automatically registers in catalog
6. Pushes to GitHub repository

## Phase 4: MCP AI-Ops Server

### Location
`manifests/bleeding-edge/mcp/`

### Features
- **Model Context Protocol (MCP)** for AI agent integration
- **GitHub Copilot** webhook integration
- **Read capabilities**: Pod logs, Hubble flows, Prometheus metrics, deployment statuses
- **Write capabilities**: Apply manifests, trigger Argo CD syncs, scale deployments, spin up vClusters
- **Granular RBAC** for safe AI operations
- **Audit logging** for all AI actions

### Tools Implemented

1. **Kubernetes Tools** (`k8s-tools.yaml`)
   - get_pod_logs
   - list_pods
   - get_deployment_status
   - apply_manifest
   - scale_deployment

2. **Hubble Tools** (`hubble-tools.yaml`)
   - get_hubble_flows
   - get_hubble_metrics

3. **Prometheus Tools** (`prometheus-tools.yaml`)
   - query_prometheus
   - get_service_metrics

4. **Argo CD Tools** (`argocd-tools.yaml`)
   - list_argocd_apps
   - get_app_status
   - trigger_sync
   - rollback_app

5. **vCluster Tools** (`vcluster-tools.yaml`)
   - list_vclusters
   - create_vcluster
   - delete_vcluster
   - connect_vcluster

### Installation
```bash
bash manifests/bleeding-edge/mcp/install-mcp-server.sh
```

### Uninstallation
```bash
bash manifests/bleeding-edge/mcp/uninstall-mcp-server.sh
```

### Usage
```bash
# Configure GitHub Copilot
# Set COPILOT_TOKEN environment variable
export COPILOT_TOKEN=your_copilot_token

# MCP Server endpoint
http://mcp-server.mcp-server.svc.cluster.local:8080/webhook

# Example GitHub Copilot prompts:
# "Show me the logs for the api-gateway pods"
# "Get the network flows from Hubble for the last 5 minutes"
# "Query Prometheus for the error rate of the api-gateway"
# "Trigger an Argo CD sync for the production application"
# "Scale the rust-data-store deployment to 3 replicas"
# "Create a vCluster for PR #123"
```

### Safety Features
- **Require confirmation** for write operations
- **Dry run mode** available
- **Audit logging** enabled
- **Rate limiting** (100 requests/minute)
- **Repository allowlist** for security

## Architecture Overview

```
+-------------------------------------------------+
|              GitHub Copilot (MCP)               |
+------------------------+------------------------+
                         | (Read/Write K8s Ops)
                         v
+-----------------------+           +--------------------+           +------------------------+
|   Backstage Portal    | --------> |   Kind Cluster     | <-------- |    Kyverno Engine      |
| (Go/Rust/Python Templates)         | (Separate Config)  |           | (Admission Validation) |
+-----------------------+           +---------+----------+           +------------------------+
                                              |
                                              v
                                  +-----------------------+
                                  | Cilium / Hubble eBPF  |
                                  | (Complementary CNI)   |
                                  +-----------------------+
```

## Installation Sequence

Install components in dependency order:

1. **Phase 1: Kyverno** (Foundation)
   ```bash
   bash manifests/bleeding-edge/kyverno/install-kyverno.sh
   ```

2. **Phase 2: Cilium + Hubble** (Infrastructure Visibility)
   ```bash
   bash manifests/bleeding-edge/cilium/install-cilium.sh
   ```

3. **Phase 3: Backstage** (Developer Self-Service)
   ```bash
   bash manifests/bleeding-edge/backstage/install-backstage.sh
   ```

4. **Phase 4: MCP Server** (Intelligent Automation)
   ```bash
   bash manifests/bleeding-edge/mcp/install-mcp-server.sh
   ```

## Access URLs

| Component | URL | Port Forward |
|-----------|-----|--------------|
| Backstage | http://localhost:3000 | `kubectl port-forward svc/backstage -n backstage 3000:80` |
| Hubble UI | http://localhost:12000 | `kubectl port-forward svc/hubble-ui -n cilium 12000:80` |
| MCP Server | Internal cluster service | N/A (webhook integration) |

## Cluster Isolation

All four components are structured in dedicated manifest packages (`manifests/bleeding-edge/`) to allow independent deployment onto the Kind environment without mutating existing state files in `k8s/base/`.

## Security Considerations

### Kyverno
- Enforces security policies at admission time
- Cosign verification ensures only signed images run
- Resource limits prevent resource exhaustion
- Rootless containers prevent privilege escalation

### Cilium
- L7 network policies for service-to-service communication
- eBPF-based security without sidecar overhead
- Network flow visibility for security monitoring

### Backstage
- Authentication and authorization for developer access
- GitHub token-based authentication
- Service account with limited Kubernetes permissions

### MCP Server
- Granular RBAC for read/write operations
- Repository allowlist for security
- Rate limiting to prevent abuse
- Audit logging for compliance
- Confirmation required for destructive operations

## Dependencies

- **Phase 1**: No dependencies (foundation layer)
- **Phase 2**: Depends on Phase 1 (security policies before network visibility)
- **Phase 3**: Depends on Phase 1 & 2 (security + observability before self-service)
- **Phase 4**: Depends on Phase 1, 2, 3 (leverages all previous layers)

## Rollback Strategy

Each phase includes uninstallation scripts for independent rollback:

```bash
# Phase 1
bash manifests/bleeding-edge/kyverno/uninstall-kyverno.sh

# Phase 2
bash manifests/bleeding-edge/cilium/uninstall-cilium.sh

# Phase 3
bash manifests/bleeding-edge/backstage/uninstall-backstage.sh

# Phase 4
bash manifests/bleeding-edge/mcp/uninstall-mcp-server.sh
```

## Integration with Existing Pipeline

### Argo CD Integration
- Backstage plugin displays Argo CD application status
- MCP Server can trigger Argo CD syncs and rollbacks
- Kyverno policies apply to Argo CD deployments

### vCluster Integration
- Backstage plugin can spin up vClusters on demand
- MCP Server can create/delete vClusters
- Kyverno policies apply to vCluster namespaces

### Prometheus Integration
- Cilium metrics exported to Prometheus
- MCP Server can query Prometheus metrics
- Backstage can display Prometheus metrics

### Chaos Mesh Integration
- Cilium network policies complement Chaos Mesh experiments
- Hubble can visualize Chaos Mesh network disruptions
- MCP Server can trigger Chaos Mesh experiments

## Troubleshooting

### Kyverno Issues
```bash
# Check policy reports
kubectl get cpol -A
kubectl get pol -A

# View policy violations
kubectl get clusterpolicyviolations
kubectl get policyviolations -A

# Check Kyverno status
kubectl get deployment -n kyverno
```

### Cilium Issues
```bash
# Check Cilium status
cilium status

# Check Hubble status
cilium hubble status

# View Cilium logs
kubectl logs -n cilium deployment/cilium-operator

# Restart Cilium
kubectl rollout restart deployment/cilium -n cilium
```

### Backstage Issues
```bash
# Check Backstage status
kubectl get deployment -n backstage

# View Backstage logs
kubectl logs -n backstage deployment/backstage

# Check database connection
kubectl get pods -n backstage
```

### MCP Server Issues
```bash
# Check MCP Server status
kubectl get deployment -n mcp-server

# View MCP Server logs
kubectl logs -n mcp-server deployment/mcp-server

# Check RBAC
kubectl get clusterrole mcp-server-role
kubectl get clusterrolebinding mcp-server-binding
```

## Performance Impact

### Kyverno
- Minimal overhead (admission webhook)
- Audit mode: ~5-10ms per request
- Enforce mode: ~10-20ms per request

### Cilium
- eBPF host routing: negligible overhead
- Hubble observability: ~1-2% CPU overhead
- Network policies: negligible latency impact

### Backstage
- Standalone deployment (no impact on workloads)
- ~500m CPU, 512Mi memory limits

### MCP Server
- Standalone deployment (no impact on workloads)
- ~500m CPU, 256Mi memory limits
- Only active during AI agent interactions

## Next Steps

1. **Configure Production Secrets**
   - Set GitHub token for Backstage
   - Set Argo CD credentials for Backstage
   - Set Copilot token for MCP Server

2. **Customize Policies**
   - Adjust Kyverno policies for your requirements
   - Switch from audit to enforce mode after validation
   - Add custom network policies in Cilium

3. **Extend Backstage Templates**
   - Add service skeleton templates
   - Customize scaffolder actions
   - Add additional plugins

4. **Configure MCP Tools**
   - Add custom MCP tools
   - Adjust rate limits
   - Configure additional safety checks

5. **Monitor and Tune**
   - Monitor Kyverno policy violations
   - Review Hubble network flows
   - Track Backstage usage metrics
   - Audit MCP Server operations

## Summary

These bleeding-edge enhancements transform the CI/CD pipeline from an enterprise-grade system to a cutting-edge platform engineering solution:

- **Kyverno**: Runtime policy-as-code admission control
- **Cilium + Hubble**: Kernel-level observability without sidecars
- **Backstage**: Developer self-service portal with golden paths
- **MCP Server**: AI-powered operations via GitHub Copilot

The pipeline now represents the absolute upper echelon of modern platform engineering, matching the capabilities of top-tier organizations while maintaining zero-cost local development through Kind.
