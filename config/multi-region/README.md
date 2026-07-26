# Multi-Region Configuration

This directory contains configuration files for multi-region disaster recovery setup and local simulation.

## Purpose

Multi-region configuration enables:
- Active-active disaster recovery
- Geographic distribution
- Local development simulation
- DR failover testing

## Contents

### Kind Cluster Configurations

- **kind-clusters/** - Kind cluster configurations for local DR simulation:
  - `primary-cluster.yaml` - Primary region simulation (us-east-1)
  - `secondary-cluster.yaml` - Secondary region simulation (eu-central-1)
  - `local-dev-cluster.yaml` - Local development cluster

### DR Simulation Scripts

Located in `scripts/dr-failover/`:
- `setup-multi-kind.sh` - Setup multi-region Kind clusters
- `simulate-failover.sh` - Simulate DR failover between regions

## Usage

### Local DR Simulation

```bash
# Setup multi-region Kind clusters
bash scripts/dr-failover/setup-multi-kind.sh

# Simulate failover
bash scripts/dr-failover/simulate-failover.sh

# Check health
bash scripts/dr-failover/simulate-failover.sh check

# Restore primary
bash scripts/dr-failover/simulate-failover.sh restore
```

### Cluster Configuration

Each Kind cluster simulates a production region:

**Primary Cluster (us-east-1)**
- 3 worker nodes
- Cilium networking
- SPIRE identity
- CoreDNS with custom zones

**Secondary Cluster (eu-central-1)**
- 3 worker nodes
- Cilium networking
- SPIRE identity
- CoreDNS with custom zones

**Local Dev Cluster**
- 1 worker node
- Basic networking
- Simplified configuration

## Multi-Region Architecture

### Production Setup

**AWS Regions:**
- Primary: us-east-1 (Virginia)
- Secondary: eu-central-1 (Frankfurt)

**Components:**
- EKS clusters in both regions
- Route53 ARC for DNS failover
- S3 cross-region replication
- SPIRE federation for identity
- Argo CD for GitOps

### Local Simulation

**Kind Clusters:**
- Simulate production topology
- Test DR procedures locally
- Validate failover logic
- No cloud costs

## DNS Configuration

### Production DNS

Route53 ARC provides:
- Health-based routing
- Automatic failover
- Latency-based routing
- Weighted routing

### Local DNS Simulation

CoreDNS configuration:
- Custom zones for simulation
- Health check endpoints
- Failover simulation
- Local testing

## Failover Process

### Automatic Failover

1. Health check detects primary failure
2. Route53 updates DNS records
3. Traffic routes to secondary region
4. Argo CD maintains cluster state
5. Services continue operating

### Manual Failover

1. Verify secondary region health
2. Update DNS configuration
3. Monitor traffic routing
4. Verify service availability
5. Document failover

### Failback Process

1. Verify primary region recovery
2. Sync data from secondary
3. Update DNS configuration
4. Monitor traffic routing
5. Verify service availability

## Best Practices

### DR Testing

- Regular failover drills
- Document procedures
- Test communication channels
- Validate recovery time objectives

### Data Consistency

- Implement data replication
- Test recovery procedures
- Monitor replication lag
- Validate data integrity

### Monitoring

- Monitor both regions
- Set up alerting
- Track failover events
- Document incidents

## Troubleshooting

### Cluster Setup Issues

Check:
- Kind installation
- Docker daemon status
- Resource availability
- Network configuration

### Failover Issues

Verify:
- DNS configuration
- Health check endpoints
- Service availability
- Network connectivity

### Data Replication Issues

Check:
- S3 bucket configuration
- Replication rules
- IAM permissions
- Network connectivity

## Documentation

- [Kind Documentation](https://kind.sigs.k8s.io/)
- [Route53 ARC Documentation](https://docs.aws.amazon.com/route53/)
- [Multi-Region Best Practices](https://aws.amazon.com/blogs/architecture/)
