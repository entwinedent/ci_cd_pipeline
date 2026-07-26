# Architecture Documentation

## System Overview

This CI/CD pipeline implements a sophisticated microservices architecture with GitOps, progressive delivery, and AI-powered monitoring. The system is designed for zero-cost local development while maintaining production-grade practices.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        GitHub Repository                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │   Source     │  │   K8s        │  │   GitHub Actions     │  │
│  │   Code       │  │   Manifests  │  │   Workflows          │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              │ GitOps Sync
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Local Kind Cluster                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │   Argo CD    │  │   Services   │  │   Argo Rollouts      │  │
│  │   Controller │  │   (Go/Rust/  │  │   (Canary Deploy)     │  │
│  │              │  │    Python)   │  │                      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Microservices Architecture

### Service Communication Flow

```
External Client
       │
       │ HTTP/1.1 (Port 8080)
       ▼
┌──────────────────┐
│  Go API Gateway  │
│  - Request       │
│    Validation    │
│  - Routing       │
│  - gRPC Client   │
└────────┬─────────┘
         │
         │ gRPC (Port 50051)
         ▼
┌──────────────────┐
│  Rust Data Store │
│  - DashMap       │
│    Cache         │
│  - TTL Management│
│  - Metrics       │
└────────┬─────────┘
         │
         │ Structured Logs
         ▼
┌──────────────────┐
│ Python Telemetry │
│  - Log Collector │
│  - Anomaly       │
│    Detection     │
│  - Webhook       │
│    Alerts        │
└──────────────────┘
```

### Service Specifications

#### Go API Gateway
- **Language**: Go 1.22+
- **Port**: 8080 (HTTP/1.1)
- **Responsibilities**:
  - External HTTP entrypoint
  - JSON request validation
  - gRPC client multiplexing
  - Health check endpoints
  - Structured logging
- **Resource Limits**: 64-128Mi RAM, 100-200m CPU
- **Image Size**: < 15MB (scratch base)

#### Rust Data Store
- **Language**: Rust 1.75+
- **Port**: 50051 (gRPC HTTP/2)
- **Responsibilities**:
  - Thread-safe in-memory storage
  - Lock-free concurrent operations (DashMap)
  - TTL expiration management
  - Time-series metrics collection
  - gRPC service implementation
- **Resource Limits**: 128-256Mi RAM, 100-500m CPU
- **Image Size**: < 25MB (debian:bookworm-slim)

#### Python Telemetry
- **Language**: Python 3.11+
- **Port**: 8000 (HTTP/JSON)
- **Responsibilities**:
  - Structured log ingestion
  - AI anomaly detection (scikit-learn)
  - Webhook-based alerting
  - Auto-remediation triggers
  - Metrics analysis
- **Resource Limits**: 256-512Mi RAM, 200-500m CPU
- **Image Size**: < 100MB (python:3.11-slim)

## Data Structures

### Rust In-Memory Store

```rust
// Primary Key-Value Store
DashMap<String, CacheEntry>  // Lock-free concurrent access

// Time-Series Ring Buffer
Arc<RwLock<VecDeque<MetricPoint>>>  // Sliding-window metrics

// TTL Expiration Map
Binary Min-Heap sorted by expiry timestamps
```

### gRPC Protocol Definition

```protobuf
service DataStoreService {
  rpc Set (SetRequest) returns (SetResponse);
  rpc Get (GetRequest) returns (GetResponse);
  rpc Delete (DeleteRequest) returns (DeleteResponse);
  rpc HealthCheck (HealthCheckRequest) returns (HealthCheckResponse);
}
```

## CI/CD Pipeline Architecture

### GitHub Actions Workflow

```
┌─────────────────────────────────────────────────────────────┐
│                    Push to GitHub                           │
└────────────────────────┬────────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
         ▼               ▼               ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│  Go Tests    │ │ Rust Tests   │ │ Python Tests │
│  - go test   │ │ - cargo test │ │ - pytest     │
│  - golangci  │ │ - clippy     │ │ - black      │
│    -lint     │ │ - fmt check  │ │ - flake8     │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Docker Build     │
              │ - Multi-platform │
              │ - Layer caching  │
              └────────┬─────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
        ▼              ▼              ▼
┌──────────────┐ ┌──────────────┐ ┌──────────────┐
│ Secret Scan  │ │ Container    │ │ SAST Scan    │
│ - Gitleaks   │ │ Scan         │ │ - Semgrep    │
│              │ │ - Trivy      │ │ - CodeQL     │
└──────┬───────┘ └──────┬───────┘ └──────┬───────┘
       │                │                │
       └────────────────┼────────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │ Push to GHCR     │
              │ - SHA tagging   │
              │ - Latest tag    │
              └────────┬─────────┘
                       │
                       ▼
              ┌──────────────────┐
              │ Update K8s       │
              │ Manifests       │
              └────────┬─────────┘
                       │
                       ▼
              ┌──────────────────┐
              │ Argo CD Sync     │
              │ - Auto-deploy    │
              └──────────────────┘
```

## Security Architecture

### Defense in Depth

1. **Code Security**
   - Secret scanning (Gitleaks)
   - Static analysis (Semgrep, CodeQL)
   - Dependency scanning (cargo audit, go mod)

2. **Container Security**
   - Vulnerability scanning (Trivy)
   - Minimal base images (scratch, distroless)
   - Rootless containers
   - Read-only filesystems

3. **Runtime Security**
   - Network policies
   - Resource limits
   - Pod security policies
   - Service mesh (future)

### Security Scanning Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│                    Security Scanning                         │
├─────────────────────────────────────────────────────────────┤
│  1. Pre-commit: Gitleaks (secret detection)                 │
│  2. CI Pipeline:                                             │
│     - Semgrep (SAST)                                         │
│     - Trivy (Container vulnerability)                        │
│     - Kube-linter (K8s configuration)                        │
│  3. Runtime:                                                 │
│     - Network policies                                       │
│     - Resource quotas                                        │
│     - Pod security standards                                 │
└─────────────────────────────────────────────────────────────┘
```

## GitOps Architecture

### Argo CD Workflow

```
┌─────────────────────────────────────────────────────────────┐
│                    Git Repository                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  k8s/base/go-gateway/deployment.yaml                 │   │
│  │  k8s/base/rust-store/deployment.yaml                  │   │
│  │  k8s/base/python-telemetry/deployment.yaml           │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ Git Push
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Argo CD Controller                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  - Detects Git changes                               │   │
│  │  - Compares desired vs actual state                   │   │
│  │  - Syncs cluster state                                │   │
│  │  - Reports sync status                                │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ kubectl apply
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                          │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  - Updates deployments                                │   │
│  │  - Rolls out new versions                             │   │
│  │  - Maintains desired state                            │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Progressive Delivery Architecture

### Canary Deployment Strategy

```
┌─────────────────────────────────────────────────────────────┐
│                    Argo Rollouts                              │
├─────────────────────────────────────────────────────────────┤
│  Step 1: Deploy 10% traffic to new version                   │
│  Step 2: Pause for analysis (30s)                            │
│  Step 3: Check metrics (success rate, latency, errors)       │
│  Step 4: If metrics OK → Scale to 50%                       │
│  Step 5: Pause for analysis (60s)                            │
│  Step 6: If metrics OK → Scale to 100%                      │
│  Step 7: If metrics fail → Auto-rollback                    │
└─────────────────────────────────────────────────────────────┘
```

### Auto-Remediation Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Anomaly Detection                          │
├─────────────────────────────────────────────────────────────┤
│  Python Telemetry detects anomaly (z-score > 3.0)           │
│  ↓                                                           │
│  Webhook trigger to Argo CD                                 │
│  ↓                                                           │
│  Argo Rollouts initiates rollback                           │
│  ↓                                                           │
│  Previous stable version restored                          │
│  ↓                                                           │
│  Alert sent to Slack/Email                                   │
└─────────────────────────────────────────────────────────────┘
```

## Monitoring & Observability

### Logging Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Structured Logging                         │
├─────────────────────────────────────────────────────────────┤
│  Go API Gateway → JSON logs → stdout → Kubernetes logs      │
│  Rust Data Store → JSON logs → stdout → Kubernetes logs      │
│  Python Telemetry → JSON logs → stdout → Kubernetes logs    │
│                          ↓                                   │
│                    Python Telemetry                           │
│                    (Ingestion & Analysis)                     │
└─────────────────────────────────────────────────────────────┘
```

### Metrics Collection

```
┌─────────────────────────────────────────────────────────────┐
│                    Metrics Flow                               │
├─────────────────────────────────────────────────────────────┤
│  Services → Metrics → Python Telemetry → Anomaly Detection   │
│                          ↓                                   │
│                    Time-series Analysis                       │
│                    (Z-score, thresholds)                      │
│                          ↓                                   │
│                    Alert Generation                          │
│                    (Webhook triggers)                         │
└─────────────────────────────────────────────────────────────┘
```

## Network Architecture

### Service Communication

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Network                         │
├─────────────────────────────────────────────────────────────┤
│  External → NodePort → go-api-gateway (8080)                │
│                          ↓                                   │
│                    ClusterIP → rust-data-store (50051)       │
│                          ↓                                   │
│                    ClusterIP → python-telemetry (8000)        │
└─────────────────────────────────────────────────────────────┘
```

### DNS Resolution

- **Service Discovery**: Kubernetes internal DNS
- **Service Names**: `go-api-gateway`, `rust-data-store`, `python-telemetry`
- **Namespace**: Default (can be configured per environment)

## Deployment Architecture

### Local Development (Kind)

```
┌─────────────────────────────────────────────────────────────┐
│                    Development Environment                    │
├─────────────────────────────────────────────────────────────┤
│  - Kind cluster (1 control-plane, 2 workers)                │
│  - Local container registry (port 5000)                    │
│  - Docker Compose for local testing                         │
│  - Port forwarding for service access                       │
│  - Zero cloud costs                                         │
└─────────────────────────────────────────────────────────────┘
```

### Production (Future)

```
┌─────────────────────────────────────────────────────────────┐
│                    Production Environment                     │
├─────────────────────────────────────────────────────────────┤
│  - Managed Kubernetes (EKS/GKE/AKS)                         │
│  - GitHub Container Registry                                │
│  - Argo CD for GitOps                                       │
│  - Argo Rollouts for progressive delivery                   │
│  - External monitoring (Prometheus/Grafana)                  │
│  - External logging (ELK/Loki)                              │
└─────────────────────────────────────────────────────────────┘
```

## Performance Considerations

### Resource Optimization

- **Go**: Static binaries, minimal runtime, < 15MB images
- **Rust**: Zero-cost abstractions, memory safety, < 25MB images
- **Python**: Multi-stage builds, dependency caching, < 100MB images

### Scalability

- **Horizontal Scaling**: Kubernetes HPA based on CPU/memory
- **Vertical Scaling**: Resource requests/limits configured
- **Database Scaling**: Rust DashMap scales with CPU cores

### Latency Targets

- **Go API Gateway**: < 100ms p95 latency
- **Rust Data Store**: < 10ms gRPC latency
- **Python Telemetry**: < 500ms anomaly detection

## Design Decisions

### Language Selection

- **Go**: Fast compilation, built-in concurrency, cloud-native ecosystem
- **Rust**: Memory safety, zero-cost abstractions, WebAssembly support
- **Python**: AI/ML libraries, rapid development, extensive tooling

### Technology Choices

- **Kind**: Local Kubernetes development, zero cost, Docker-native
- **Argo CD**: Declarative GitOps, Kubernetes-native, excellent UI
- **GitHub Actions**: Free tier, tight Git integration, extensive marketplace
- **gRPC**: Type-safe contracts, high performance, streaming support

### Architecture Patterns

- **GitOps**: Single source of truth, automated sync, audit trail
- **Microservices**: Language polyglot, independent scaling, fault isolation
- **Progressive Delivery**: Zero-downtime deployments, automated rollback
- **Shift-Left Security**: Early detection, automated blocking, compliance

## Future Enhancements

### Phase 9-10 Roadmap

- [ ] Argo Rollouts integration
- [ ] Advanced anomaly detection models
- [ ] Service mesh (Istio/Linkerd)
- [ ] Distributed tracing (Jaeger/Tempo)
- [ ] External monitoring integration
- [ ] Multi-environment support
- [ ] Helm charts for all services
- [ ] Performance benchmarking suite

---

This architecture demonstrates production-grade DevOps practices while maintaining zero-cost local development capabilities.
