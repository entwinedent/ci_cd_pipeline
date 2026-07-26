# ZAP Security Scanning Configuration

This directory contains OWASP ZAP (Zed Attack Proxy) configuration files for automated security scanning of web applications.

## Purpose

ZAP is used to perform automated security scans on the CI/CD pipeline web services to identify common vulnerabilities:
- SQL injection
- Cross-site scripting (XSS)
- Security misconfigurations
- Information disclosure
- Authentication issues

## Contents

### Rules Files

- **rules.tsv** - General security scanning rules configuration
- **rules-telemetry.tsv** - Specific rules for telemetry service scanning

## Usage

### Local Scanning

```bash
# Start ZAP daemon
zap.sh -daemon -port 8080

# Run automated scan
zap-cli quick-scan --self-contained http://localhost:8080

# Generate report
zap-cli report -o zap-report.html
```

### CI/CD Integration

ZAP scans are integrated into the security scanning workflow:

```yaml
- name: Run ZAP Security Scan
  run: |
    zap-cli start
    zap-cli quick-scan --self-contained http://localhost:8080
    zap-cli report -o zap-report.html
```

## Configuration

### Scan Rules

Configure which security checks to perform in the rules files:

- **Active Scan Rules** - Aggressive vulnerability detection
- **Passive Scan Rules** - Non-intrusive analysis
- **Authentication Rules** - Login/logout testing
- **Session Management** - Cookie and session handling

### Target Configuration

Configure scan targets:
- API Gateway: `http://localhost:8080`
- Telemetry: `http://localhost:8000`
- Data Store: `http://localhost:50051`

## Best Practices

### Scan Frequency

- Run on every pull request
- Scan before deployment
- Regular scheduled scans
- After major code changes

### False Positives

- Document expected false positives
- Configure exclusions in rules files
- Review and update regularly
- Maintain security exception process

### Performance

- Use passive scans for quick checks
- Active scans for comprehensive testing
- Limit scan duration in CI
- Cache scan results when possible

## Troubleshooting

### Scan Failures

Common issues:
- Target service not running
- Network connectivity issues
- Authentication failures
- Timeout errors

### False Positives

Handle false positives by:
- Reviewing scan results
- Configuring exclusions
- Documenting exceptions
- Updating rules files

## Documentation

- [OWASP ZAP Documentation](https://www.zaproxy.org/docs/)
- [ZAP CLI Documentation](https://github.com/zaproxy/zaproxy)
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
