# GitHub Configuration

This directory contains GitHub-specific configuration files for the CI/CD pipeline repository.

## Contents

### Workflows (`.github/workflows/`)

Contains GitHub Actions workflow definitions for CI/CD automation:

- **ci.yml** - Main CI pipeline with testing, linting, and Docker builds
- **security-scan.yml** - Security scanning with Gitleaks, Trivy, Semgrep, and Snyk
- **deploy.yml** - Deployment pipeline with image signing and registry push
- **finops-gate.yml** - Cost analysis and FinOps governance
- **phylum-gate.yml** - Supply chain security analysis with Phylum

### Dependabot (`.github/dependabot.yml`)

Automated dependency update configuration:
- GitHub Actions updates
- Go module updates
- Rust crate updates
- Python package updates
- Docker base image updates

## Usage

### CI/CD Pipeline

Workflows automatically trigger on:
- Push to main branch
- Pull requests
- Manual workflow dispatch

### Security Scanning

Security scans run on:
- Every pull request
- Push to main branch
- Manual trigger

### Deployment

Deployment occurs when:
- Tests pass
- Security scans complete
- Code is merged to main branch

## Configuration

### Secrets Required

Configure these secrets in GitHub repository settings:

- `GHCR_TOKEN` - GitHub Container Registry token
- `SNYK_TOKEN` - Snyk API token
- `PHYLUM_API_KEY` - Phylum API key
- `COSIGN_PRIVATE_KEY` - Cosign signing key
- `SLACK_WEBHOOK_URL` - Slack notifications

### Environment Variables

Configure these environment variables in workflow files:

- `DOCKER_REGISTRY` - Container registry URL
- `IMAGE_TAG` - Docker image tag
- `DEPLOY_ENVIRONMENT` - Deployment environment

## Best Practices

### Workflow Optimization

- Use matrix strategies for parallel execution
- Cache dependencies to speed up builds
- Use reusable actions for common tasks
- Implement proper error handling

### Security

- Pin action versions to commit SHAs
- Use least privilege for secrets
- Enable branch protection rules
- Require status checks for merges

### Maintenance

- Regularly update dependencies
- Monitor workflow execution times
- Review security scan results
- Update action versions

## Troubleshooting

### Workflow Failures

Check workflow logs for:
- Test failures
- Build errors
- Security scan issues
- Deployment problems

### Secret Issues

Verify secrets are:
- Properly configured in repository settings
- Not expired
- Have correct permissions
- Are accessible to workflows

### Dependency Updates

Review Dependabot PRs for:
- Breaking changes
- Security vulnerabilities
- License compatibility
- Test compatibility

## Documentation

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Dependabot Documentation](https://docs.github.com/en/code-security/dependabot)
- [GitHub Security Features](https://docs.github.com/en/code-security)
