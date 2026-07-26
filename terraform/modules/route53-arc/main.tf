# Route53 Application Recovery Controller Module
# Provides global DNS routing with health-based failover

resource "aws_route53_zone" "this" {
  name = var.domain_name
}

resource "aws_route53_record" "api_primary" {
  zone_id = aws_route53_zone.this.zone_id
  name    = "api"
  type    = "CNAME"
  ttl     = "60"
  records = [var.primary_lb_dns]
}

resource "aws_route53_record" "api_secondary" {
  zone_id = aws_route53_zone.this.zone_id
  name    = "api-secondary"
  type    = "CNAME"
  ttl     = "60"
  records = [var.secondary_lb_dns]
}

resource "aws_route53_health_check" "primary" {
  fqdn              = var.primary_health_check_endpoint
  port              = 443
  type              = "HTTPS"
  resource_path     = "/healthz"
  request_interval  = 30
  failure_threshold = 3
}

resource "aws_route53_health_check" "secondary" {
  fqdn              = var.secondary_health_check_endpoint
  port              = 443
  type              = "HTTPS"
  resource_path     = "/healthz"
  request_interval  = 30
  failure_threshold = 3
}

resource "aws_route53_record" "failover" {
  zone_id = aws_route53_zone.this.zone_id
  name    = "api-failover"
  type    = "CNAME"
  ttl     = "30"
  
  failover_routing_policy {
    type = "PRIMARY"
  }
  
  set_identifier = "primary"
  alias {
    name                   = var.primary_lb_dns
    zone_id                = var.primary_lb_zone_id
    evaluate_target_health = true
  }
  
  health_check_id = aws_route53_health_check.primary.id
}

resource "aws_route53_record" "failover_secondary" {
  zone_id = aws_route53_zone.this.zone_id
  name    = "api-failover"
  type    = "CNAME"
  ttl     = "30"
  
  failover_routing_policy {
    type = "SECONDARY"
  }
  
  set_identifier = "secondary"
  alias {
    name                   = var.secondary_lb_dns
    zone_id                = var.secondary_lb_zone_id
    evaluate_target_health = true
  }
  
  health_check_id = aws_route53_health_check.secondary.id
}

# Variables
variable "domain_name" {
  description = "Domain name for Route53 zone"
  type        = string
}

variable "primary_lb_dns" {
  description = "Primary load balancer DNS name"
  type        = string
}

variable "secondary_lb_dns" {
  description = "Secondary load balancer DNS name"
  type        = string
}

variable "primary_lb_zone_id" {
  description = "Primary load balancer zone ID"
  type        = string
}

variable "secondary_lb_zone_id" {
  description = "Secondary load balancer zone ID"
  type        = string
}

variable "primary_health_check_endpoint" {
  description = "Primary health check endpoint"
  type        = string
}

variable "secondary_health_check_endpoint" {
  description = "Secondary health check endpoint"
  type        = string
}

# Outputs
output "zone_id" {
  description = "Route53 hosted zone ID"
  value       = aws_route53_zone.this.zone_id
}

output "zone_name" {
  description = "Route53 hosted zone name"
  value       = aws_route53_zone.this.name
}

output "api_dns" {
  description = "API DNS name"
  value       = aws_route53_record.api_primary.fqdn
}

output "failover_dns" {
  description = "Failover DNS name"
  value       = aws_route53_record.failover.fqdn
}
