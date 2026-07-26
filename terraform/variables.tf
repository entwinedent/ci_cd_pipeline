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

variable "vpc_cidr_primary" {
  description = "CIDR block for primary VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "vpc_cidr_secondary" {
  description = "CIDR block for secondary VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "cluster_name_primary" {
  description = "Name for primary EKS cluster"
  type        = string
  default     = "ci-cd-pipeline-primary"
}

variable "cluster_name_secondary" {
  description = "Name for secondary EKS cluster"
  type        = string
  default     = "ci-cd-pipeline-secondary"
}
