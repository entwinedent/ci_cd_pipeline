# Observability & Runtime Identity Screenshots

## Required Screenshots

### 1. Grafana / Tempo Distributed Trace (`grafana-trace.png`)
**What to capture:**
- Tempo distributed trace showing end-to-end request flow
- Trace starting at go-api-gateway
- gRPC call to rust-data-store
- HTTP call to python-telemetry
- Service-to-service timing and latency information
- W3C trace context propagation

**How to capture:**
1. Generate test traffic through all services
2. Access Grafana/Tempo UI (typically `http://localhost:3000` for Grafana)
3. Navigate to Tempo traces
4. Find a trace that spans all three services
5. Capture the trace detail view

**Expected content:**
- Trace timeline showing service hops
- Service names and spans visible
- Latency information for each hop
- Trace ID and parent/child relationships
- Error indicators (if any)

### 2. Hubble UI Network Topology (`hubble-network-topology.png`)
**What to capture:**
- Hubble UI network topology map
- Real-time L7 flows between services
- Service connections (go-api-gateway ↔ rust-data-store ↔ python-telemetry)
- Active Cilium network policies
- Flow direction and protocol information

**How to capture:**
1. Ensure Cilium and Hubble are installed
2. Access Hubble UI (typically `http://localhost:8082`)
3. Navigate to Network Graph or Flow view
4. Generate traffic between services
5. Capture the network topology showing active flows

**Expected content:**
- Service nodes with connections
- Flow indicators showing direction
- Protocol information (HTTP, gRPC)
- Network policy indicators
- Real-time flow statistics

### 3. SPIFFE/SPIRE Identity Status (`spire-identity-status.png`)
**What to capture:**
- Terminal output from `spire-server entry show` command
- Issued SVID cryptographic identities
- Active mTLS bindings
- SPIFFE IDs for each service
- Trust bundle information

**How to capture:**
1. Access the cluster where SPIRE is installed
2. Run `kubectl exec -n spire spire-server-0 -- ./spire-server entry show`
3. Capture the terminal output showing identity information
4. Optionally show SVID details for specific services

**Expected content:**
- SPIFFE IDs for each service
- SVID certificate information
- Trust bundle details
- Registration status
- Selector information
