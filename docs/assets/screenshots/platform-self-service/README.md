# Platform Self-Service & Control Screenshots

## Required Screenshots

### 1. Backstage Service Catalog (`backstage-catalog.png`)
**What to capture:**
- Backstage service catalog main interface
- Cataloged microservices (go-api-gateway, rust-data-store, python-telemetry)
- Service metadata and health status
- "Golden Path" templates for new services
- Service ownership and documentation links

**How to capture:**
1. Access Backstage UI (typically `http://localhost:3000`)
2. Navigate to the Catalog page
3. Ensure all microservices are registered
4. Capture the main catalog view showing all services

**Expected content:**
- Service cards for each microservice
- Service health indicators
- Service metadata (owner, description, tags)
- Template section visible
- Search/filter interface

### 2. Unleash Feature Flag Dashboard (`unleash-feature-flags.png`)
**What to capture:**
- Unleash feature flag dashboard
- Active feature flags (rate_limiting, experimental_cache)
- Toggle states (enabled/disabled)
- Feature flag strategies (gradual rollout, canary, etc.)
- Provider toggle states (Unleash vs LaunchDarkly integration)

**How to capture:**
1. Access Unleash dashboard (typically `http://localhost:4242`)
2. Navigate to the Features page
3. Ensure feature flags are configured
4. Capture the feature flags overview

**Expected content:**
- Feature flag list with toggle states
- Feature flag descriptions and strategies
- Environment-specific configurations
- Usage statistics (if available)
- Integration indicators (LaunchDarkly, etc.)
