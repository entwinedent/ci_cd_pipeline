# Main Terraform configuration for Multi-Region Active-Active DR
# This module deploys EKS clusters in us-east-1 and eu-central-1 with
# cross-region networking, Route53 ARC, and disaster recovery infrastructure

terraform {
  required_version = ">= 1.5.0"
  
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.12"
    }
  }

  backend "s3" {
    bucket         = "ci-cd-pipeline-terraform-state"
    key            = "multi-region-dr/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "ci-cd-pipeline-terraform-locks"
  }
}

provider "aws" {
  region = var.primary_region
  
  default_tags {
    tags = {
      Project     = "ci-cd-pipeline"
      Environment = var.environment
      ManagedBy   = "terraform"
      DRStrategy  = "active-active"
    }
  }
}

provider "aws" {
  alias  = "secondary"
  region = var.secondary_region
  
  default_tags {
    tags = {
      Project     = "ci-cd-pipeline"
      Environment = var.environment
      ManagedBy   = "terraform"
      DRStrategy  = "active-active"
    }
  }
}

# Variables
variable "primary_region" {
  description = "Primary AWS region"
  type        = string
  default     = "us-east-1"
}

variable "secondary_region" {
  description = "Secondary AWS region"
  type        = string
  default     = "eu-central-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "cluster_version" {
  description = "EKS cluster version"
  type        = string
  default     = "1.28"
}

variable "domain_name" {
  description = "Domain name for Route53"
  type        = string
  default     = "ci-cd-pipeline.local"
}

# Outputs
output "primary_cluster_endpoint" {
  description = "Primary EKS cluster endpoint"
  value       = module.eks_primary.cluster_endpoint
}

output "secondary_cluster_endpoint" {
  description = "Secondary EKS cluster endpoint"
  value       = module.eks_secondary.cluster_endpoint
}

output "primary_cluster_name" {
  description = "Primary EKS cluster name"
  value       = module.eks_primary.cluster_name
}

output "secondary_cluster_name" {
  description = "Secondary EKS cluster name"
  value       = module.eks_secondary.cluster_name
}

output "route53_zone_id" {
  description = "Route53 hosted zone ID"
  value       = module.route53.zone_id
}
