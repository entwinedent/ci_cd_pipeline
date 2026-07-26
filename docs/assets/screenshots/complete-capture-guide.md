# Complete Screenshot Capture Guide

## Prerequisites Check

Currently Missing:
- ❌ Kind (Kubernetes in Docker)
- ❌ GitHub CLI (gh)

These tools are required to capture the platform UI screenshots and trigger GitHub Actions workflows.

## Installation Instructions

### Install Kind
**Windows (Chocolatey):**
```powershell
choco install kind
```

**Windows (Manual):**
1. Download from https://kind.sigs.k8s.io/
2. Add to PATH
3. Verify: `kind version`

### Install GitHub CLI
**Windows (Chocolatey):**
```powershell
choco install gh
```

**Windows (Manual):**
1. Download from https://cli.github.com/
2. Add to PATH
3. Authenticate: `gh auth login`

## Step-by-Step Screenshot Capture Process

### Phase 1: GitHub Actions Screenshots (No Kind Required)

#### 1. Trigger CI/CD Pipeline
```bash
# Create a test branch
git checkout -b screenshot-capture

# Push to trigger workflows
git push origin screenshot-capture

# Create PR
gh pr create --title "Screenshot Capture" --body "Triggering workflows for screenshot capture"
```

#### 2. Capture CI/CD Matrix Screenshot
1. Go to GitHub Actions tab
2. Click on "CI Pipeline" workflow
3. Wait for all jobs to complete
4. Screenshot the jobs overview showing green checkmarks
5. Save as: `screenshots/ci-cd-security/github-actions-matrix.png`

#### 3. Capture Security Scan Results
1. Click on "Security Scan" workflow
2. Navigate through individual security jobs
3. Capture the security scan summary
4. Save as: `screenshots/ci-cd-security/security-scan-results.png`

#### 4. Capture FinOps PR Comment
1. Wait for FinOps Cost Governance Gate to complete
2. Go to PR conversation tab
3. Find the automated cost analysis comment
4. Capture the comment
5. Save as: `screenshots/ci-cd-security/finops-pr-comment.png`

### Phase 2: Platform Services Deployment (Requires Kind)

#### 1. Create Kind Cluster
```bash
kind create cluster --config config/kind-config.yaml
```

#### 2. Deploy Cilium + Hubble
```bash
bash manifests/bleeding-edge/cilium/install-cilium.sh
```

#### 3. Deploy SPIRE
```bash
bash manifests/bleeding-edge/spire/install-spire.sh
```

#### 4. Deploy Argo CD
```bash
bash scripts/install-argocd.sh
```

#### 5. Deploy Backstage
```bash
bash manifests/bleeding-edge/backstage/install-backstage.sh
```

#### 6. Deploy Unleash
```bash
kubectl apply -f manifests/bleeding-edge/unleash/unleash-deployment.yaml
kubectl apply -f manifests/bleeding-edge/unleash/unleash-config.yaml
```

### Phase 3: Platform UI Screenshots

#### 1. Argo CD Dashboard
```bash
# Port forward Argo CD
kubectl port-forward svc/argocd-server -n argocd 8080:443
```
1. Open http://localhost:8080
2. Login with admin credentials
3. Capture the dashboard showing synced applications
4. Save as: `screenshots/gitops-progressive-delivery/argocd-dashboard.png`

#### 2. Argo Rollouts Canary
1. Deploy a canary rollout using existing manifests
2. Access Argo Rollouts dashboard
3. Capture during active canary deployment
4. Save as: `screenshots/gitops-progressive-delivery/argo-rollouts-canary.png`

#### 3. Hubble Network Topology
```bash
# Port forward Hubble UI
kubectl port-forward svc/hubble-ui -n cilium 12000:80
```
1. Open http://localhost:12000
2. Generate traffic between services
3. Capture network topology
4. Save as: `screenshots/observability-runtime-identity/hubble-network-topology.png`

#### 4. SPIRE Identity Status
```bash
# Show SPIRE entries
kubectl exec -n spire spire-server-0 -- ./spire-server entry show
```
1. Capture terminal output showing SVID identities
2. Save as: `screenshots/observability-runtime-identity/spire-identity-status.png`

#### 5. Backstage Catalog
```bash
# Port forward Backstage
kubectl port-forward svc/backstage -n backstage 3000:80
```
1. Open http://localhost:3000
2. Capture service catalog
3. Save as: `screenshots/platform-self-service/backstage-catalog.png`

#### 6. Unleash Feature Flags
```bash
# Port forward Unleash
kubectl port-forward svc/unleash -n default 4242:4242
```
1. Open http://localhost:4242
2. Capture feature flag dashboard
3. Save as: `screenshots/platform-self-service/unleash-feature-flags.png`

### Phase 4: Update README.md

Once all screenshots are captured, they will automatically appear in the README.md since the image references are already in place.

## Alternative: Use Existing Workflow Runs

If you have existing GitHub Actions runs:
1. Go to Actions tab
2. Click on recent successful runs
3. Capture screenshots from those runs
4. Place in appropriate directories

## Current Limitations

- Kind installed via Docker
- GitHub CLI not available for triggering workflows
- Platform services can be deployed via Kind cluster

## Next Steps

1. Install Kind and GitHub CLI
2. Follow this guide step-by-step
3. Capture screenshots as specified
4. Verify screenshots appear in README.md
