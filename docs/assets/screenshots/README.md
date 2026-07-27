# Portfolio Screenshots

This directory contains essential screenshots for demonstrating the CI/CD pipeline capabilities across four operational domains.

## Current Status

- ✅ Directory structure created
- ✅ Documentation for each screenshot category
- ✅ Kind installed via Docker
- ✅ Installation scripts available for all platform services
- ✅ Docker Compose services running locally

## Deployment Prerequisites

To capture platform UI screenshots, you need:
- Kind (Kubernetes in Docker)
- Helm 3.x
- kubectl
- Docker or Docker Desktop

See `deployment-guide.md` for installation instructions and deployment scripts.

## Screenshot Categories

### A. CI/CD & Security Pipelines (GitHub Actions)

**Required Screenshots:**

1. **Green CI/CD Pipeline Matrix** (`images/resume screenshot 1.png`) ✅
   - Screenshot of a passing PR build showing parallel jobs
   - Include: Go/Rust/Python tests, Playwright API/E2E, Pact contracts, k6 performance, OWASP ZAP, Phylum supply chain, Cosign signing
   - Capture: GitHub Actions workflow run page with all 23 jobs passing
   - **Status:** Captured and integrated into main README

2. **Supply Chain Security Report** (`images/resume screenshot 3.png`) ✅
   - Visual outputs of Gitleaks, Trivy container CVE scans, Snyk dependency checks
   - Include: Cosign image signing verification
   - Capture: Security scan workflow results page
   - **Status:** Captured and integrated into main README

3. **Pipeline Execution Detail** (`images/resume screenshot 2.png`) ✅
   - Detailed view of contract testing, dynamic security scanning, and load testing
   - Include: Pact contract tests, Playwright E2E, k6 performance across Go, Rust, Python
   - Capture: GitHub Actions job execution detail view
   - **Status:** Captured and integrated into main README

4. **FinOps Cost Diff Comment** (`ci-cd-security/finops-pr-comment.png`) ⏳
   - Automated PR comment from Infracost/Kubecost showing predicted monthly cost deltas
   - Capture: GitHub PR comment with cost analysis before merging
   - **Status:** Not yet captured

### B. GitOps & Progressive Delivery

**Required Screenshots:**

1. **Argo CD Dashboard** (`gitops-progressive-delivery/argocd-dashboard.png`)
   - Sync status showing healthy, synced state across multi-region environments
   - Include: Application sets and cluster status
   - Capture: Argo CD UI main dashboard

2. **Argo Rollouts Canary Analysis** (`gitops-progressive-delivery/argo-rollouts-canary.png`)
   - Real-time traffic shifting (20% → 50% → 100%)
   - Include: Prometheus error-rate analysis templates during progressive deployment
   - Capture: Argo Rollouts analysis dashboard during canary deployment

### C. Observability & Runtime Identity

**Required Screenshots:**

1. **Grafana / Tempo Distributed Trace** (`observability-runtime-identity/grafana-trace.png`)
   - Single W3C trace starting at go-api-gateway, hopping via gRPC to rust-data-store, pushing to python-telemetry
   - Capture: Tempo trace view showing service-to-service calls

2. **Hubble UI (Cilium eBPF)** (`observability-runtime-identity/hubble-network-topology.png`)
   - Real-time network topology map visualizing zero-overhead L7 flows
   - Include: Active Cilium network policies between services
   - Capture: Hubble UI network topology view

3. **SPIFFE/SPIRE Identity Status** (`observability-runtime-identity/spire-identity-status.png`)
   - CLI output (spire-server entry show) demonstrating issued SVID cryptographic identities
   - Include: Active mTLS bindings
   - Capture: Terminal output showing SPIRE identity information

### D. Platform Self-Service & Control

**Required Screenshots:**

1. **Backstage Service Catalog & Scaffolder** (`platform-self-service/backstage-catalog.png`)
   - Developer portal interface showing cataloged microservices
   - Include: "Golden Path" templates for new Go, Rust, and Python services
   - Capture: Backstage UI service catalog page

2. **Unleash / Feature Flag Dashboard** (`platform-self-service/unleash-feature-flags.png`)
   - Active feature flags (rate_limiting, experimental_cache) showing provider toggle states
   - Include: Toggle states between Unleash and LaunchDarkly
   - Capture: Unleash feature flag dashboard

## Naming Conventions

- Use lowercase with hyphens: `category-description.png`
- Keep descriptions concise but descriptive
- Use PNG format for better quality
- Recommended resolution: 1920x1080 or higher

## GitHub README Integration

### Expandable HTML Blocks

```markdown
<details>
<summary>Click to view CI/CD Pipeline Matrix</summary>

![CI/CD Pipeline Matrix](screenshots/ci-cd-security/github-actions-matrix.png)

</details>
```

### 2-Column Image Grid

```markdown
<table>
  <tr>
    <td><b>Argo CD Dashboard</b><br><img src="screenshots/gitops-progressive-delivery/argocd-dashboard.png" width="400"></td>
    <td><b>Argo Rollouts Canary</b><br><img src="screenshots/gitops-progressive-delivery/argo-rollouts-canary.png" width="400"></td>
  </tr>
</table>
```

## Capture Instructions

### GitHub Actions Screenshots
1. Trigger a PR build with all workflows
2. Navigate to Actions tab in GitHub
3. Capture the workflow run page showing all passing jobs
4. For security scans, capture the detailed results page
5. For FinOps, capture the PR comment with cost analysis

### Argo CD Screenshots
1. Access Argo CD dashboard (typically `http://localhost:8080` for local)
2. Navigate to Applications page
3. Capture the dashboard showing sync status
4. For Rollouts, capture during an active canary deployment

### Observability Screenshots
1. Access Grafana/Tempo UI
2. Generate a test trace through all services
3. Capture the distributed trace view
4. For Hubble, capture the network topology during active traffic
5. For SPIRE, run `spire-server entry show` and capture terminal output

### Platform Screenshots
1. Access Backstage service catalog
2. Capture the main catalog page with registered services
3. For Unleash, access feature flag dashboard
4. Capture the active flags and their states

## Placeholder Files

Each subdirectory contains a placeholder file with specific capture instructions for that category. Refer to these files when capturing screenshots.
