# GitOps & Progressive Delivery Screenshots

## Required Screenshots

### 1. Argo CD Dashboard (`argocd-dashboard.png`)
**What to capture:**
- Argo CD main dashboard showing application sync status
- Multi-region cluster status (primary and secondary)
- Application sets with health indicators
- Sync status showing "Synced" and "Healthy" states

**How to capture:**
1. Access Argo CD dashboard (typically `http://localhost:8080` for local Kind clusters)
2. Navigate to Applications or Dashboard view
3. Ensure all applications show healthy status
4. Capture the full dashboard view

**Expected content:**
- Green health indicators for all applications
- Sync status showing "Synced"
- Cluster information visible
- Application groups/sets organized

### 2. Argo Rollouts Canary Analysis (`argo-rollouts-canary.png`)
**What to capture:**
- Argo Rollouts analysis dashboard during canary deployment
- Traffic shifting visualization (20% → 50% → 100%)
- Prometheus error-rate analysis metrics
- Success rate and latency metrics
- Canary vs stable version comparison

**How to capture:**
1. Trigger a canary deployment (via Argo Rollouts or kubectl apply)
2. Access Argo Rollouts dashboard
3. Navigate to the specific rollout's analysis page
4. Capture during active canary deployment showing traffic shift

**Expected content:**
- Traffic percentage bars (20%, 50%, 100%)
- Error rate graphs showing stable vs canary
- Success rate comparison
- Analysis metrics (latency, throughput)
- Rollout status and progress
