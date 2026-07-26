#!/bin/bash

# DR Failover Simulation Script
# Simulates failover from primary to secondary region

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

# Check cluster health
check_cluster_health() {
    local cluster=$1
    kubectl config use-context "kind-$cluster"
    
    print_info "Checking health of $cluster..."
    
    # Check control plane
    if kubectl get nodes &> /dev/null; then
        print_success "$cluster control plane is healthy"
    else
        print_error "$cluster control plane is unhealthy"
        return 1
    fi
    
    # Check pods
    if kubectl get pods -A &> /dev/null; then
        print_success "$cluster pods are running"
    else
        print_error "$cluster pods are not running"
        return 1
    fi
}

# Simulate primary region failure
simulate_primary_failure() {
    print_info "Simulating primary region failure..."
    kubectl config use-context kind-ci-cd-pipeline-primary
    
    # Scale down deployments to simulate failure
    kubectl scale deployment --all --replicas=0 -A
    print_success "Primary region simulated failure"
}

# Failover to secondary region
failover_to_secondary() {
    print_info "Failing over to secondary region..."
    kubectl config use-context kind-ci-cd-pipeline-secondary
    
    # Scale up deployments
    kubectl scale deployment --all --replicas=2 -A
    print_success "Failover to secondary region complete"
}

# Restore primary region
restore_primary() {
    print_info "Restoring primary region..."
    kubectl config use-context kind-ci-cd-pipeline-primary
    
    # Scale up deployments
    kubectl scale deployment --all --replicas=2 -A
    print_success "Primary region restored"
}

# Main script
case "${1:-check}" in
    check)
        print_info "Checking cluster health..."
        check_cluster_health "ci-cd-pipeline-primary"
        check_cluster_health "ci-cd-pipeline-secondary"
        check_cluster_health "ci-cd-pipeline-local"
        ;;
    failover)
        print_info "Starting DR failover simulation..."
        check_cluster_health "ci-cd-pipeline-primary"
        simulate_primary_failure
        failover_to_secondary
        print_success "DR failover simulation complete"
        ;;
    restore)
        print_info "Restoring primary region..."
        restore_primary
        print_success "Primary region restoration complete"
        ;;
    *)
        echo "Usage: $0 {check|failover|restore}"
        echo "  check    - Check health of all clusters"
        echo "  failover - Simulate primary failure and failover to secondary"
        echo "  restore  - Restore primary region"
        exit 1
        ;;
esac
