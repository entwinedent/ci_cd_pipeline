# GitHub Actions Screenshot Capture Guide

## Manual Capture Instructions

Since GitHub CLI (gh) is not accessible in the current shell path, follow these manual steps to capture GitHub Actions screenshots:

### Step 1: Access GitHub Actions
1. Open your web browser
2. Navigate to your GitHub repository
3. Click on the "Actions" tab in the repository navigation

### Step 2: Select Workflow Run
1. Choose a recent successful workflow run
2. Look for workflows like:
   - "CI Pipeline" or "Build and Test"
   - "Security Scan"
   - "FinOps Cost Governance Gate"
3. Click on the workflow run to view details

### Step 3: Capture Pipeline DAG Screenshot
1. Scroll to see the workflow visualization
2. Capture the entire pipeline DAG showing:
   - All jobs with their status (green checkmarks)
   - Job dependencies and parallel execution
   - Execution times
3. Save as: `screenshots/ci-cd-security/github-actions-matrix.png`

### Step 4: Capture Security Scan Results
1. Navigate to the "Security Scan" workflow run
2. Click on individual security jobs to see detailed results
3. Capture screenshots showing:
   - Gitleaks secret scan results
   - Trivy container vulnerability scan
   - Snyk dependency vulnerability report
   - Cosign image signing verification
4. Save as: `screenshots/ci-cd-security/security-scan-results.png`

### Step 5: Capture FinOps PR Comment (if available)
1. Navigate to a PR that triggered the FinOps workflow
2. Go to the PR conversation tab
3. Find the automated cost analysis comment
4. Capture the comment showing:
   - Cost breakdown (infrastructure vs Kubernetes)
   - Monthly cost delta
   - Percentage increase
   - Approval status
5. Save as: `screenshots/ci-cd-security/finops-pr-comment.png`

## Alternative: Use Existing Workflow Runs

If you don't have recent workflow runs, you can:
1. Create a test branch: `git checkout -b screenshot-test`
2. Make a small change and push: `git push origin screenshot-test`
3. Create a PR: `gh pr create --title "Screenshot Test" --body "Testing workflow for screenshots"`
4. Wait for workflows to complete
5. Capture screenshots from the new workflow runs

## Screenshot Tips

- Use high-resolution capture (1920x1080 or higher)
- Ensure all text is readable
- Include browser window for context
- Capture full workflow visualization when possible
- Save in PNG format for best quality

## Expected Content

### GitHub Actions Matrix
- Green checkmarks for all jobs
- Job names: go-test, rust-test, python-test, playwright-e2e, pact-contracts, k6-performance
- Matrix strategy configuration visible
- Execution times for each job

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
