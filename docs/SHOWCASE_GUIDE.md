# Portfolio Showcase Guide

This guide provides instructions for creating screenshots and demo assets to showcase the CI/CD pipeline portfolio.

## Required Screenshots

### 1. Architecture Diagram
- **File**: `docs/images/architecture-diagram.png`
- **Content**: High-level system architecture showing all components
- **Creation**: Use the Mermaid diagram from README.md or create a visual diagram
- **Tools**: Mermaid Live Editor, Draw.io, or Lucidchart

### 2. Grafana OpenTelemetry Dashboard
- **File**: `docs/images/grafana-dashboard.png`
- **Content**: Distributed tracing dashboard with service spans
- **Creation**: 
  ```bash
  # Deploy Tempo and Grafana
  kubectl port-forward svc/grafana 3000:3000
  # Navigate to Tempo dashboard
  # Take screenshot of trace visualization
  ```

### 3. FinOps Cost Dashboard
- **File**: `docs/images/finops-dashboard.png`
- **Content**: OpenCost cost breakdown by namespace/service
- **Creation**:
  ```bash
  # Deploy OpenCost
  kubectl port-forward svc/opencost 9003:9003
  # Navigate to OpenCost UI
  # Take screenshot of cost dashboard
  ```

### 4. Argo Rollouts Canary Deployment
- **File**: `docs/images/argo-rollouts.png`
- **Content**: Canary deployment visualization with traffic splitting
- **Creation**:
  ```bash
  # Deploy Argo Rollouts
  kubectl port-forward svc/argo-rollouts-dashboard 3100:3100
  # Navigate to Rollouts dashboard
  # Take screenshot of canary deployment
  ```

### 5. Backstage Developer Portal
- **File**: `docs/images/backstage-portal.png`
- **Content**: Developer portal showing microservice templates
- **Creation**:
  ```bash
  # Deploy Backstage
  kubectl port-forward svc/backstage 7007:7007
  # Navigate to Backstage UI
  # Take screenshot of catalog and templates
  ```

### 6. Chaos Engineering Demo
- **File**: `docs/images/chaos-testing.gif`
- **Content**: Terminal GIF showing chaos experiment execution
- **Creation**:
  ```bash
  # Use terminal recording tools
  # Linux: asciinema or ttyrec
  # macOS: terminalizer or recordit
  # Windows: ShareX or OBS
  ```

### 7. k6 Load Test Results
- **File**: `docs/images/k6-load-test.png`
- **Content**: Load test metrics and performance graphs
- **Creation**:
  ```bash
  # Run k6 load test
  k6 run tests/load-test.js
  # Take screenshot of results
  ```

## Demo Script

### Terminal Recording Script

Create a comprehensive demo showing:

```bash
# 1. Show repository structure
tree -L 2 -I 'node_modules|__pycache__|target'

# 2. Run tests
make test-unit

# 3. Build and deploy
make docker-build
make kind-setup
make deploy

# 4. Check status
make status

# 5. Port forward and test
make port-forward
curl http://localhost:8080/healthz
curl http://localhost:8080/api/v1/data/test

# 6. Show Argo CD
kubectl port-forward svc/argocd-server -n argocd 8080:443
# Show Argo CD UI in browser

# 7. Show distributed tracing
kubectl port-forward svc/tempo 16686:16686
# Show Tempo UI in browser

# 8. Run chaos test
make test-chaos
```

## Image Guidelines

### Technical Specifications
- **Format**: PNG for static images, GIF for animations
- **Size**: Keep under 500KB for web optimization
- **Dimensions**: 1920x1080 (16:9 aspect ratio)
- **Quality**: High resolution for readability

### Content Guidelines
- **Clarity**: Ensure text is readable at 100% zoom
- **Context**: Include relevant UI elements and navigation
- **Consistency**: Use same theme/color scheme across screenshots
- **Annotations**: Add arrows/callouts to highlight important features

### Privacy Considerations
- **Redact**: Remove any sensitive information (API keys, tokens)
- **Sanitize**: Use mock data instead of real credentials
- **Anonymize**: Remove personal information from logs

## Creation Tools

### Screenshot Tools
- **Linux**: `gnome-screenshot`, `shutter`, `flameshot`
- **macOS**: Built-in screenshot (Cmd+Shift+4), CleanShot X
- **Windows**: Snipping Tool, ShareX, Greenshot

### Terminal Recording
- **Linux**: `asciinema`, `ttyrec`, `terminalizer`
- **macOS**: `terminalizer`, `recordit`, `iTerm2 recording`
- **Windows**: `ShareX`, `OBS Studio`, `PowerToys`

### Diagram Tools
- **Mermaid**: [Mermaid Live Editor](https://mermaid.live/)
- **Draw.io**: [draw.io](https://app.diagrams.net/)
- **Lucidchart**: [lucidchart.com](https://www.lucidchart.com/)

## Automated Screenshots

### Using Playwright

```javascript
// screenshots.js
const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  // Grafana Dashboard
  await page.goto('http://localhost:3000');
  await page.screenshot({ path: 'docs/images/grafana-dashboard.png' });
  
  // OpenCost Dashboard
  await page.goto('http://localhost:9003');
  await page.screenshot({ path: 'docs/images/finops-dashboard.png' });
  
  await browser.close();
})();
```

### Using Selenium

```python
# screenshots.py
from selenium import webdriver
from selenium.webdriver.common.by import By

driver = webdriver.Chrome()
driver.get('http://localhost:3000')
driver.save_screenshot('docs/images/grafana-dashboard.png')
driver.quit()
```

## README Integration

Update the main README to reference screenshots:

```markdown
## 📸 Showcase

![Architecture Diagram](docs/images/architecture-diagram.png)

![Grafana Dashboard](docs/images/grafana-dashboard.png)

![FinOps Dashboard](docs/images/finops-dashboard.png)

## 🎥 Demo

[![Chaos Testing Demo](docs/images/chaos-testing.gif)]
```

## Checklist

- [ ] Create architecture diagram
- [ ] Capture Grafana dashboard screenshot
- [ ] Capture FinOps dashboard screenshot
- [ ] Capture Argo Rollouts screenshot
- [ ] Capture Backstage portal screenshot
- [ ] Record chaos testing demo
- [ ] Capture k6 load test results
- [ ] Optimize image sizes
- [ ] Update README with screenshots
- [ ] Verify all images load correctly
