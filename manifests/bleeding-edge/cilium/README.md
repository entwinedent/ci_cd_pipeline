# Cilium eBPF Networking

Cilium with eBPF for Kubernetes networking, security, and observability with Hubble for distributed tracing.

## Architecture

### eBPF Technology

Cilium uses eBPF (Extended Berkeley Packet Filter) for:
- Kernel-space network policy enforcement
- L7-aware traffic filtering
- High-performance network observability
- Zero-overhead load balancing

### Components

**Cilium Agent**
- Runs as DaemonSet on each node
- Manages eBPF programs in kernel
- Enforces network policies
- Provides network observability

**Hubble Relay**
- Centralized observability server
- Aggregates metrics from all nodes
- Provides distributed tracing
- Real-time network topology

**Hubble UI**
- Web-based visualization
- Service dependency graph
- Network flow monitoring
- Security policy debugging

## Deployment

### Install Cilium

```bash
cilium install --version 1.14.0
cilium status
```

### Enable Hubble

```bash
cilium hubble enable
cilium hubble ui
```

### Verify Installation

```bash
# Check Cilium status
cilium status

# Check Hubble status
kubectl get pods -n kube-system -l k8s-app=hubble-relay

# Access Hubble UI
cilium hubble ui
```

## Network Policies

### Default Policies

Cilium enforces zero-trust networking:
- Default deny all traffic
- Explicit allow rules for service communication
- L7-aware HTTP/gRPC policies

### Example Policy

```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: api-gateway-policy
spec:
  endpointSelector:
    matchLabels:
      app: go-api-gateway
  ingress:
  - fromEndpoints:
    - matchLabels:
        k8s:io.kubernetes.pod.namespace: default
    toPorts:
    - ports:
      - port: "8080"
        protocol: TCP
      rules:
        http:
          - method: GET
            path: "/api/v1/data/*"
```

## Observability

### Hubble CLI

```bash
# Observe network flows
hubble observe

# Filter by service
hubble observe --pod go-api-gateway

# Filter by protocol
hubble observe --protocol http

# View metrics
hubble metric latency
hubble metric drop
```

### Hubble UI

Access the Hubble UI:
```bash
cilium hubble ui
```

Features:
- Service dependency graph
- Real-time flow visualization
- Network policy debugging
- DNS resolution tracking

## Performance

### eBPF Benefits

**Kernel-Space Processing**
- No context switches between kernel and user space
- Direct packet processing in kernel
- Minimal overhead compared to iptables

**L7-Aware Policies**
- HTTP/gRPC protocol awareness
- Header-based filtering
- Request/response inspection

**Observability**
- Built-in metrics collection
- Distributed tracing integration
- Network topology discovery

### Benchmarks

- **Throughput**: > 10M PPS per node
- **Latency**: < 100μs policy enforcement
- **CPU Overhead**: < 5% per node
- **Memory**: < 100MB per node

## Security

### Zero-Trust Networking

- Default deny all traffic
- Explicit allow rules only
- Service-to-service mTLS
- Network policy as code

### Integration with SPIRE

Cilium integrates with SPIRE for identity:
- SPIFFE IDs in network policies
- mTLS certificate validation
- Identity-based access control

## Troubleshooting

### Connectivity Issues

```bash
# Check Cilium status
cilium status

# Check network policies
cilium network policy list

# Check connectivity
cilium connectivity test
```

### Policy Not Applied

```bash
# Check policy status
cilium network policy get <policy-name>

# Check endpoint status
cilium endpoint list

# View policy logs
cilium log warning --filter 'policy'
```

## Documentation

- [Cilium Documentation](https://docs.cilium.io/)
- [Hubble Documentation](https://docs.cilium.io/en/stable/observability/hubble/)
- [eBPF Technology](https://ebpf.io/)
