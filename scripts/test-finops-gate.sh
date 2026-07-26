#!/bin/bash

# FinOps Gate Test Script
# Simulates the cost estimation bot for high-compute resource changes

set -e

echo "=== FinOps Gate Testing ==="

# 1. Analyze current resource costs
echo "1. Analyzing Kubernetes resource costs..."
python3 << 'EOF'
import json

# Hardcoded values for testing (simulating YAML parsing)
current_monthly_cost = 100
current_cpu = "500m"
current_memory = "512Mi"

print(f"Current monthly cost: ${current_monthly_cost}")
print(f"Current CPU request: {current_cpu}")
print(f"Current memory request: {current_memory}")

# High-compute values
high_cpu = "4"
high_memory = "8Gi"
high_monthly_cost = 2500

print(f"High-compute monthly cost: ${high_monthly_cost}")
print(f"High CPU request: {high_cpu}")
print(f"High memory request: {high_memory}")

# Calculate cost delta
cost_delta = high_monthly_cost - current_monthly_cost
percentage_increase = (cost_delta / current_monthly_cost) * 100

print(f"\nCost delta: ${cost_delta}")
print(f"Percentage increase: {percentage_increase:.2f}%")

# Check thresholds
dollar_threshold = 2000
percentage_threshold = 15

requires_approval = cost_delta > dollar_threshold or percentage_increase > percentage_threshold

print(f"\nDollar threshold: ${dollar_threshold}")
print(f"Percentage threshold: {percentage_threshold}%")
print(f"Requires approval: {requires_approval}")

# Generate mock cost report
report = {
    'summary': {
        'currentMonthlyCost': current_monthly_cost,
        'newMonthlyCost': high_monthly_cost,
        'costDelta': cost_delta,
        'percentageIncrease': percentage_increase,
        'currency': 'USD'
    },
    'thresholds': {
        'percentageThreshold': percentage_threshold,
        'dollarThreshold': dollar_threshold,
        'requiresApproval': requires_approval
    }
}

with open('/tmp/finops-test-report.json', 'w') as f:
    json.dump(report, f, indent=2)

print(f"\nCost report saved to /tmp/finops-test-report.json")
print(f"Report: {json.dumps(report, indent=2)}")
EOF

# 2. Display cost impact report
echo -e "\n2. Cost Impact Report:"
echo "💰 FinOps Cost Impact Report"
echo ""
echo "### Summary"
echo "- Current Monthly Cost: $100.00"
echo "- New Monthly Cost: $2,500.00"
echo "- Cost Delta: $2,400.00"
echo "- Percentage Increase: 2400.00%"
echo ""
echo "### Cost Analysis"
echo "⚠️ This PR exceeds cost thresholds and requires FinOps-Admins approval"
echo ""
echo "### Thresholds"
echo "- Percentage Threshold: 15%"
echo "- Dollar Threshold: $2,000/month"
echo ""
echo "### Recommendations"
echo "- Review resource allocations in Kubernetes manifests"
echo "- Consider using spot instances for non-critical workloads"
echo "- Implement autoscaling to optimize resource utilization"

echo -e "\n=== FinOps Gate Testing Complete ==="
