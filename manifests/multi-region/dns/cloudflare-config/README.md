# Cloudflare Magic WAN Configuration

This directory contains the Cloudflare Magic WAN configuration for global DNS routing and traffic management.

## Architecture

Cloudflare Magic WAN provides:
- **Anycast DNS**: Global DNS routing with automatic failover
- **Health Probes**: Layer 7 health checks for cluster endpoints
- **Traffic Steering**: Intelligent traffic routing based on health and latency
- **DDoS Protection**: Built-in protection against DDoS attacks

## Configuration

### DNS Records

```
api.ci-cd-pipeline.local -> Primary Region (us-east-1)
api-secondary.ci-cd-pipeline.local -> Secondary Region (eu-central-1)
api-failover.ci-cd-pipeline.local -> Automatic failover endpoint
```

### Health Checks

Health checks monitor:
- HTTP/HTTPS endpoints on port 443
- `/healthz` endpoint for liveness
- Response time thresholds (< 500ms)
- HTTP status codes (200-299)

### Failover Logic

1. **Normal Operation**: Traffic routed to primary region
2. **Primary Degraded**: Health check fails → traffic routed to secondary
3. **Primary Restored**: Health check passes → traffic routed back to primary

## Local Simulation

For local development without Cloudflare, use CoreDNS configuration:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns
  namespace: kube-system
data:
  Corefile: |
    .:53 {
        errors
        health
        rewrite name api.ci-cd-pipeline.local go-api-gateway.default.svc.cluster.local
        rewrite name api-secondary.ci-cd-pipeline.local go-api-gateway-secondary.default.svc.cluster.local
        kubernetes cluster.local in-addr.arpa ip6.arpa {
            pods insecure
            fallthrough in-addr.arpa ip6.arpa
        }
        prometheus :9153
        forward . /etc/resolv.conf
        cache 30
        loop
        reload
        loadbalance
    }
```

## Cilium eBPF Health Probes

Combine Cloudflare health checks with Cilium Hubble metrics:

```bash
# Monitor service health with Hubble
hubble observe --protocol http --follow

# Check latency metrics
hubble metric latency

# Monitor packet drops
hubble metric drop
```

## Integration with Route53 ARC

For production, integrate with AWS Route53 ARC:

```yaml
resource "aws_route53_health_check" "primary" {
  fqdn              = "api.ci-cd-pipeline.local"
  port              = 443
  type              = "HTTPS"
  resource_path     = "/healthz"
  request_interval  = 30
  failure_threshold = 3
}
```

## Testing

### Local Testing

```bash
# Test DNS resolution
nslookup api.ci-cd-pipeline.local

# Test health endpoint
curl -k https://api.ci-cd-pipeline.local/healthz

# Test failover
kubectl scale deployment --all --replicas=0 -A  # Simulate failure
curl -k https://api-failover.ci-cd-pipeline.local/healthz  # Should route to secondary
```

### Cloudflare Testing

```bash
# Test DNS propagation
dig api.ci-cd-pipeline.local

# Test health checks
curl -I https://api.ci-cd-pipeline.local/healthz

# Monitor health check status
cloudflare-cli health-checks list
```

## Security

- **TLS**: All health checks use HTTPS
- **Authentication**: Health endpoints require valid SPIRE SVID certificates
- **Rate Limiting**: Health check endpoints rate-limited to prevent abuse
- **IP Whitelisting**: Only Cloudflare IPs allowed to access health endpoints
