# CI/CD & Security Screenshots

## Required Screenshots

### 1. GitHub Actions Matrix (`resume screenshot 1.png`)
**What to capture:**
- GitHub Actions workflow run page showing parallel jobs
- All jobs passing (green checkmarks)
- Job names visible: Go Tests, Rust Tests, Python Tests, Playwright E2E, Pact Contracts, k6 Performance, OWASP ZAP, Phylum Supply Chain, Cosign Signing

**How to capture:**
1. Create a PR or push to trigger CI/CD pipeline
2. Navigate to Actions tab in GitHub
3. Click on the CI Pipeline workflow run
4. Take screenshot of the jobs overview page

**Expected content:**
- Green checkmarks for all 23 jobs
- Job execution times visible
- Matrix strategy visible (if applicable)

**Status:** ✅ Captured - Located in `../images/resume screenshot 1.png`

### 2. Security Scan Results (`resume screenshot 3.png`)
**What to capture:**
- Security scanning workflow results
- Gitleaks secret scan output
- Trivy container CVE scan results
- Snyk dependency vulnerability reports
- Cosign image signing verification

**How to capture:**
1. Trigger security-scan.yml workflow
2. Navigate to the workflow run
3. Click on individual security jobs to see detailed results
4. Capture the security scan summary page

**Expected content:**
- Zero critical vulnerabilities (ideal)
- Security scan summary with pass/fail status
- SBOM generation confirmation

**Status:** ✅ Captured - Located in `../images/resume screenshot 3.png`

### 3. Pipeline Execution Detail (`resume screenshot 2.png`)
**What to capture:**
- Detailed pipeline execution showing contract testing
- Dynamic security scanning across Go, Rust, Python
- Load testing and performance metrics
- Supply chain security gates

**How to capture:**
1. Navigate to CI Pipeline workflow run
2. Click on individual job logs
3. Capture the execution detail view

**Expected content:**
- Contract test results (Pact)
- E2E test results (Playwright)
- Load test results (k6)
- Security scan results

**Status:** ✅ Captured - Located in `../images/resume screenshot 2.png`

### 4. FinOps PR Comment (`finops-pr-comment.png`)
**What to capture:**
- Automated PR comment from FinOps gate
- Cost delta analysis showing monthly cost changes
- Threshold warnings if costs exceed limits
- Approval requirements

**How to capture:**
1. Create a PR with resource changes
2. Wait for FinOps gate workflow to complete
3. Navigate to the PR conversation
4. Capture the automated cost analysis comment

**Expected content:**
- Cost breakdown (infrastructure vs Kubernetes)
- Monthly cost delta
- Percentage increase
- Approval status (requires approval or within thresholds)

**Status:** ⏳ Pending - Not yet captured
