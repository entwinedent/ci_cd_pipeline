variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
}

variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
}

variable "cluster_name" {
  description = "Cluster name for tagging"
  type        = string
}

variable "peer_vpc_id" {
  description = "Peer VPC ID for peering"
  type        = string
  default     = null
}

variable "peer_vpc_cidr" {
  description = "Peer VPC CIDR block"
  type        = string
  default     = null
}

variable "peer_region" {
  description = "Peer region for peering"
  type        = string
  default     = null
}
