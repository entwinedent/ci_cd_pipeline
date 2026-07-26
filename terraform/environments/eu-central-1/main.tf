# Secondary Region (eu-central-1) EKS Cluster Configuration

module "vpc_secondary" {
  source = "../../modules/vpc-networking"

  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
  cluster_name       = var.cluster_name
}

module "eks_secondary" {
  source = "../../modules/eks-cluster"

  cluster_name           = var.cluster_name
  cluster_version        = var.cluster_version
  vpc_id                 = module.vpc_secondary.vpc_id
  private_subnet_ids     = module.vpc_secondary.private_subnet_ids
  environment            = var.environment
  allowed_cidr_blocks    = var.allowed_cidr_blocks
  node_group_desired_size = var.node_group_desired_size
  node_group_max_size    = var.node_group_max_size
  node_group_min_size    = var.node_group_min_size
  instance_types         = var.instance_types
}

# Variables
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "availability_zones" {
  description = "Availability zones"
  type        = list(string)
  default     = ["eu-central-1a", "eu-central-1b", "eu-central-1c"]
}

variable "cluster_name" {
  description = "Cluster name"
  type        = string
  default     = "ci-cd-pipeline-secondary"
}

variable "cluster_version" {
  description = "EKS version"
  type        = string
  default     = "1.28"
}

variable "environment" {
  description = "Environment"
  type        = string
  default     = "production"
}

variable "allowed_cidr_blocks" {
  description = "Allowed CIDR blocks"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "node_group_desired_size" {
  description = "Desired node group size"
  type        = number
  default     = 2
}

variable "node_group_max_size" {
  description = "Maximum node group size"
  type        = number
  default     = 4
}

variable "node_group_min_size" {
  description = "Minimum node group size"
  type        = number
  default     = 1
}

variable "instance_types" {
  description = "Instance types"
  type        = list(string)
  default     = ["t3.medium"]
}

# Outputs
output "cluster_endpoint" {
  value = module.eks_secondary.cluster_endpoint
}

output "cluster_name" {
  value = module.eks_secondary.cluster_name
}

output "vpc_id" {
  value = module.vpc_secondary.vpc_id
}
