#!/usr/bin/env python3
"""
Mock Infracost Implementation for Terraform Cost Analysis
Provides cost estimation without requiring actual Infracost API
"""

import json
import sys
import argparse
from typing import Dict, List, Any
from dataclasses import dataclass
from pathlib import Path


@dataclass
class ResourceCost:
    """Cost information for a single resource"""
    resource_type: str
    resource_name: str
    monthly_cost: float
    hourly_cost: float
    region: str


@dataclass
class CostBreakdown:
    """Complete cost breakdown for Terraform configuration"""
    total_monthly_cost: float
    total_hourly_cost: float
    resources: List[ResourceCost]
    currency: str = "USD"


class MockPricingEngine:
    """Mock pricing engine for AWS resources"""
    
    # Mock pricing data (USD per hour)
    PRICING_DATA = {
        "aws_eks_cluster": {
            "us-east-1": 0.10,
            "eu-central-1": 0.12,
            "default": 0.10
        },
        "aws_eks_node_group": {
            "t3.medium": 0.0416,
            "t3.large": 0.0832,
            "m5.large": 0.096,
            "default": 0.05
        },
        "aws_nat_gateway": {
            "us-east-1": 0.045,
            "eu-central-1": 0.045,
            "default": 0.045
        },
        "aws_eip": {
            "us-east-1": 0.005,
            "eu-central-1": 0.005,
            "default": 0.005
        },
        "aws_lb": {
            "us-east-1": 0.0225,
            "eu-central-1": 0.025,
            "default": 0.0225
        },
        "aws_db_instance": {
            "t3.micro": 0.0116,
            "t3.small": 0.023,
            "t3.medium": 0.046,
            "default": 0.03
        },
        "aws_s3_bucket": {
            "storage": 0.023,  # per GB
            "requests": 0.0004,  # per 1000 requests
            "default": 0.01
        }
    }
    
    @classmethod
    def get_price(cls, resource_type: str, instance_type: str = None, region: str = "us-east-1") -> float:
        """Get hourly price for a resource type"""
        pricing_dict = cls.PRICING_DATA.get(resource_type, {})
        
        if instance_type and isinstance(pricing_dict, dict):
            price = pricing_dict.get(instance_type, pricing_dict.get("default", 0.01))
        else:
            price = pricing_dict.get(region, pricing_dict.get("default", 0.01))
        
        return price


def parse_terraform_plan(plan_file: str) -> List[Dict[str, Any]]:
    """Parse Terraform plan JSON file"""
    try:
        with open(plan_file, 'r') as f:
            plan_data = json.load(f)
        
        resources = []
        
        # Extract resource changes from plan
        if 'resource_changes' in plan_data:
            for change in plan_data['resource_changes']:
                resource = {
                    'type': change['type'],
                    'name': change['name'],
                    'address': change['address'],
                    'mode': change.get('mode', 'managed'),
                    'change': change.get('change', {})
                }
                resources.append(resource)
        
        return resources
    
    except FileNotFoundError:
        print(f"Error: Plan file not found: {plan_file}")
        sys.exit(1)
    except json.JSONDecodeError:
        print(f"Error: Invalid JSON in plan file: {plan_file}")
        sys.exit(1)


def estimate_costs(resources: List[Dict[str, Any]], region: str = "us-east-1") -> CostBreakdown:
    """Estimate costs for Terraform resources"""
    cost_resources = []
    total_hourly = 0.0
    
    for resource in resources:
        resource_type = resource['type']
        resource_name = resource['name']
        
        # Get base price
        hourly_cost = MockPricingEngine.get_price(resource_type, region=region)
        
        # Adjust for resource count if specified
        change = resource.get('change', {})
        if 'after' in change:
            after = change['after']
            if isinstance(after, dict):
                count = after.get('count', 1)
                if isinstance(count, int):
                    hourly_cost *= count
        
        # Add some randomness for realism
        hourly_cost *= (0.9 + (hash(resource_name) % 20) / 100.0)
        
        monthly_cost = hourly_cost * 730  # 730 hours per month
        
        cost_resources.append(ResourceCost(
            resource_type=resource_type,
            resource_name=resource_name,
            monthly_cost=monthly_cost,
            hourly_cost=hourly_cost,
            region=region
        ))
        
        total_hourly += hourly_cost
    
    total_monthly = total_hourly * 730
    
    return CostBreakdown(
        total_monthly_cost=total_monthly,
        total_hourly_cost=total_hourly,
        resources=cost_resources,
        currency="USD"
    )


def format_currency(amount: float) -> str:
    """Format amount as currency"""
    return f"${amount:.2f}"


def generate_report(cost_breakdown: CostBreakdown, output_format: str = "json") -> str:
    """Generate cost report in specified format"""
    if output_format == "json":
        report = {
            "totalMonthlyCost": cost_breakdown.total_monthly_cost,
            "totalHourlyCost": cost_breakdown.total_hourly_cost,
            "currency": cost_breakdown.currency,
            "resources": [
                {
                    "type": r.resource_type,
                    "name": r.resource_name,
                    "monthlyCost": r.monthly_cost,
                    "hourlyCost": r.hourly_cost,
                    "region": r.region
                }
                for r in cost_breakdown.resources
            ]
        }
        return json.dumps(report, indent=2)
    
    elif output_format == "text":
        lines = [
            "=== Terraform Cost Breakdown ===",
            f"Total Monthly Cost: {format_currency(cost_breakdown.total_monthly_cost)}",
            f"Total Hourly Cost: {format_currency(cost_breakdown.total_hourly_cost)}",
            f"Currency: {cost_breakdown.currency}",
            "",
            "Resource Breakdown:"
        ]
        
        for resource in cost_breakdown.resources:
            lines.append(f"  - {resource.resource_type}.{resource.resource_name}")
            lines.append(f"    Monthly: {format_currency(resource.monthly_cost)}")
            lines.append(f"    Hourly: {format_currency(resource.hourly_cost)}")
            lines.append(f"    Region: {resource.region}")
        
        return "\n".join(lines)
    
    else:
        raise ValueError(f"Unknown output format: {output_format}")


def main():
    parser = argparse.ArgumentParser(description="Mock Infracost for Terraform cost analysis")
    parser.add_argument("--path", required=True, help="Path to Terraform plan JSON file")
    parser.add_argument("--format", default="json", choices=["json", "text"], help="Output format")
    parser.add_argument("--region", default="us-east-1", help="AWS region for pricing")
    parser.add_argument("--out-file", help="Output file path")
    
    args = parser.parse_args()
    
    # Parse Terraform plan
    resources = parse_terraform_plan(args.path)
    
    # Estimate costs
    cost_breakdown = estimate_costs(resources, region=args.region)
    
    # Generate report
    report = generate_report(cost_breakdown, output_format=args.format)
    
    # Output report
    if args.out_file:
        with open(args.out_file, 'w') as f:
            f.write(report)
        print(f"Cost report written to {args.out_file}")
    else:
        print(report)


if __name__ == "__main__":
    main()
