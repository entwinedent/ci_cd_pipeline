# GitHub Actions Screenshots - Capture Instructions

Since GitHub Actions runs in the cloud, you can capture these screenshots without needing Kind or local Kubernetes.

## How to Capture GitHub Actions Screenshots

### 1. Trigger CI/CD Pipeline Matrix
1. Create a new PR or push to trigger the CI/CD workflows
2. Navigate to the Actions tab in your GitHub repository
3. Click on the "CI Pipeline" workflow run
4. Scroll to see all the parallel jobs (Go Tests, Rust Tests, Python Tests, Playwright E2E, Pact Contracts, k6 Performance)
5. Take a screenshot showing all jobs with green checkmarks
6. Save as: `ci-cd-security/github-actions-matrix.png`

### 2. Capture Security Scan Results
1. Navigate to the "Security Scan" workflow run
2. Click on individual security jobs to see detailed results
3. For Gitleaks: Capture the secret scan output
4. For Trivy: Capture the container vulnerability scan results
5. For Snyk: Capture the dependency vulnerability report
6. For Cosign: Capture the image signing verification
7. Save as: `ci-cd-security/security-scan-results.png`

### 3. Capture FinOps PR Comment
1. Create a PR with resource changes (e.g., modify `k8s/base/high-compute-test.yaml`)
2. Wait for the FinOps Cost Governance Gate workflow to complete
3. Navigate to the PR conversation tab
4. Find the automated cost analysis comment from the FinOps bot
5. Capture the comment showing cost delta and approval status
6. Save as: `ci-cd-security/finops-pr-comment.png`

## Alternative: Use Existing Workflow Runs

If you have existing workflow runs, you can capture screenshots from those:
1. Go to Actions tab
2. Click on recent successful workflow runs
3. Capture the relevant job results and comments

## Expected Content

### GitHub Actions Matrix
- Green checkmarks for all jobs
- Job names visible: go-test, rust-test, python-test, playwright-e2e, pact-contracts, k6-performance
- Execution times visible
- Matrix strategy configuration visible

### Security Scan Results
- Zero critical vulnerabilities (ideal)
- Security scan summary with pass/fail status
- SBOM generation confirmation
- Image signing verification success

### FinOps PR Comment
- Cost breakdown (infrastructure vs Kubernetes)
- Monthly cost delta (e.g., $2,400 increase)
- Percentage increase (e.g., 2400%)
- Approval status (requires approval or within thresholds)
- Recommendations for cost optimization
