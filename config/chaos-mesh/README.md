# Chaos Mesh Configuration

This directory contains Chaos Mesh configuration files for chaos engineering experiments to test system resilience.

## Purpose

Chaos Mesh is used to inject faults into the CI/CD pipeline services to:
- Test system resilience
- Identify failure modes
- Validate recovery mechanisms
- Improve system reliability

## Contents

### Installation

- **install-chaos-mesh.sh** - Script to install Chaos Mesh in Kubernetes cluster

### Chaos Experiments

- **pod-kill-api-gateway.yaml** - Kill API Gateway pods
- **pod-kill-rust-store.yaml** - Kill Data Store pods
- **memory-stress-api-gateway.yaml** - Memory stress on API Gateway
- **memory-stress-rust-store.yaml** - Memory stress on Data Store
- **network-delay.yaml** - Network delay injection
- **network-delay-rust-store.yaml** - Network delay on Data Store

## Usage

### Installation

```bash
# Install Chaos Mesh
bash config/chaos-mesh/install-chaos-mesh.sh

# Verify installation
kubectl get pods -n chaos-mesh
```

### Running Experiments

```bash
# Apply chaos experiment
kubectl apply -f config/chaos-mesh/pod-kill-api-gateway.yaml

# Check experiment status
kubectl get chaosexperiments

# Delete experiment
kubectl delete chaosexperiment pod-kill-api-gateway
```

### Chaos Dashboard

Access the Chaos Mesh dashboard:

```bash
kubectl port-forward svc/chaos-dashboard -n chaos-mesh 2333:2333
```

## Experiment Types

### Pod Faults

- **Pod Kill** - Terminates pods to test pod recovery
- **Pod Failure** - Simulates pod failures
- **Container Kill** - Terminates specific containers

### Network Faults

- **Network Delay** - Adds latency to network traffic
- **Network Loss** - Drops network packets
- **Network Partition** - Isolates network segments

### Resource Faults

- **Memory Stress** - Consumes memory resources
- **CPU Stress** - Consumes CPU resources
- **Disk Stress** - Consumes disk I/O

### Application Faults

- **HTTP Faults** - Injects HTTP errors
- **DNS Faults** - Manipulates DNS resolution
- **JVM Faults** - Java-specific faults (not applicable)

## Best Practices

### Experiment Design

- Start with simple experiments
- Gradually increase complexity
- Test one fault at a time
- Document expected outcomes

### Safety Measures

- Use namespace isolation
- Set appropriate duration limits
- Monitor system health during experiments
- Have rollback procedures ready

### CI/CD Integration

- Run chaos tests in staging environment
- Integrate with monitoring and alerting
- Automate experiment execution
- Collect metrics and analysis

## Configuration

### Experiment Parameters

Configure experiment parameters:
- **Duration** - How long the fault lasts
- **Selector** - Which pods to target
- **Severity** - Impact level of the fault
- **Schedule** - When to run experiments

### Safety Limits

Set safety limits:
- Maximum experiment duration
- Target pod limits
- Namespace restrictions
- Resource usage limits

## Troubleshooting

### Experiment Not Applying

Check:
- Chaos Mesh installation status
- Experiment YAML syntax
- Target pod availability
- RBAC permissions

### System Not Recovering

Verify:
- Health check configuration
- Pod restart policies
- Resource limits
- Monitoring alerts

### Dashboard Issues

Check:
- Dashboard service status
- Port forwarding configuration
- Browser console errors
- Network connectivity

## Documentation

- [Chaos Mesh Documentation](https://chaos-mesh.org/docs/)
- [Chaos Engineering Best Practices](https://principlesofchaos.org/)
- [Kubernetes Chaos Engineering](https://kubernetes.io/docs/concepts/cluster-administration/chaos-engineering/)
