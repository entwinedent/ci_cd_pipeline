# Terraform Infrastructure

Multi-region AWS infrastructure for production-grade disaster recovery with EKS clusters, VPC networking, and Route53 health-based routing.

## Architecture

### Multi-Region Setup

**Primary Region (us-east-1)**
- EKS cluster with managed node groups
- VPC with public/private subnets
- NAT gateways for private subnet egress
- Cross-region VPC peering

**Secondary Region (eu-central-1)**
- EKS cluster with managed node groups
- VPC with public/private subnets
- NAT gateways for private subnet egress
- Cross-region VPC peering

### Components

**VPC Module** (`modules/vpc-networking/`)
- Creates VPC with CIDR block
- Public and private subnets across AZs
- Internet gateway and NAT gateways
- Route tables for subnet routing
- Cross-region VPC peering

**EKS Module** (`modules/eks-cluster/`)
- EKS cluster with specified version
- Managed node groups with auto-scaling
- IAM roles for cluster and nodes
- Security groups for cluster access
- Cluster logging enabled

**Route53 ARC Module** (`modules/route53-arc/`)
- Route53 hosted zone for domain
- Health checks for cluster endpoints
- Failover routing policies
- DNS records for API endpoints

**S3 Replication Module** (`modules/s3-replication/`)
- Primary and secondary S3 buckets
- Cross-region replication configuration
- Versioning and encryption enabled
- Lifecycle policies for cost optimization

## Usage

### Initialize Terraform

```bash
cd terraform
terraform init
```

### Plan Changes

```bash
terraform plan -out=tfplan
```

### Apply Changes

```bash
terraform apply tfplan
```

### Destroy Infrastructure

```bash
terraform destroy
```

## Configuration

### Variables

Edit `terraform.tfvars` or pass variables:

```hcl
primary_region   = "us-east-1"
secondary_region = "eu-central-1"
cluster_version  = "1.28"
domain_name      = "ci-cd-pipeline.local"
environment      = "production"
```

### Outputs

After deployment, Terraform outputs:
- Cluster endpoints
- Cluster names
- Route53 zone ID
- S3 bucket ARNs

## Cost Estimation

Use the mock Infracost tool:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > /tmp/plan.json
python ../tools/infracost-mock/infracost_mock.py --path /tmp/plan.json
```

## Security

- IAM roles with least privilege
- VPC private subnets for workloads
- Security groups with minimal rules
- S3 bucket encryption enabled
- EKS control plane private access

## Multi-Region DR

### Failover Process

1. Route53 health checks detect primary failure
2. DNS automatically routes to secondary region
3. Argo CD maintains cluster state in both regions
4. SPIRE federation maintains trust across regions

### Testing DR

Use the Kind cluster simulation:

```bash
bash scripts/dr-failover/setup-multi-kind.sh
bash scripts/dr-failover/simulate-failover.sh
```

## Dependencies

- Terraform >= 1.5.0
- AWS CLI configured with credentials
- kubectl for cluster access
