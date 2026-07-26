#!/usr/bin/env python3
"""
Mock kubectl-cost Implementation for Kubernetes Manifest Cost Prediction
Provides cost estimation for Kubernetes workloads without requiring actual Kubecost
"""

import json
import sys
import argparse
import yaml
from typing import Dict, List, Any
from dataclasses import dataclass
from pathlib import Path


@dataclass
class WorkloadCost:
    """Cost information for a Kubernetes workload"""
    workload_type: str
    workload_name: str
    namespace: str
    monthly_cost: float
    cpu_cost: float
    memory_cost: float
    storage_cost: float


@dataclass
class CostPrediction:
    """Complete cost prediction for Kubernetes manifests"""
    total_monthly_cost: float
    total_cpu_cost: float
    total_memory_cost: float
    total_storage_cost: float
    workloads: List[WorkloadCost]
    currency: str = "USD"


class MockKubernetesPricing:
    """Mock pricing for Kubernetes resources"""
    
    # Pricing per unit (USD per month)
    CPU_PRICE_PER_CORE = 30.0  # $30 per vCPU core per month
    MEMORY_PRICE_PER_GB = 4.0  # $4 per GB RAM per month
    STORAGE_PRICE_PER_GB = 0.10  # $0.10 per GB storage per month
    
    # Instance type mappings
    INSTANCE_SPECS = {
        "t3.micro": {"cpu": 2, "memory": 1},  # 2 vCPU, 1GB
        "t3.small": {"cpu": 2, "memory": 2},   # 2 vCPU, 2GB
        "t3.medium": {"cpu": 2, "memory": 4},  # 2 vCPU, 4GB
        "t3.large": {"cpu": 2, "memory": 8},   # 2 vCPU, 8GB
        "m5.large": {"cpu": 2, "memory": 8},    # 2 vCPU, 8GB
        "m5.xlarge": {"cpu": 4, "memory": 16},  # 4 vCPU, 16GB
    }
    
    @classmethod
    def calculate_resource_cost(cls, cpu_cores: float, memory_gb: float, storage_gb: float = 0) -> Dict[str, float]:
        """Calculate cost for Kubernetes resources"""
        cpu_cost = cpu_cores * cls.CPU_PRICE_PER_CORE
        memory_cost = memory_gb * cls.MEMORY_PRICE_PER_GB
        storage_cost = storage_gb * cls.STORAGE_PRICE_PER_GB
        
        return {
            "cpu_cost": cpu_cost,
            "memory_cost": memory_cost,
            "storage_cost": storage_cost,
            "total_cost": cpu_cost + memory_cost + storage_cost
        }


def parse_kubernetes_manifests(manifest_path: str) -> List[Dict[str, Any]]:
    """Parse Kubernetes YAML manifests"""
    manifests = []
    
    try:
        path = Path(manifest_path)
        
        if path.is_file():
            files = [path]
        elif path.is_dir():
            files = list(path.glob("*.yaml")) + list(path.glob("*.yml"))
        else:
            print(f"Error: Invalid path: {manifest_path}")
            sys.exit(1)
        
        for file in files:
            try:
                with open(file, 'r') as f:
                    docs = list(yaml.safe_load_all(f))
                    
                    for doc in docs:
                        if doc and isinstance(doc, dict):
                            manifests.append(doc)
            
            except yaml.YAMLError as e:
                print(f"Warning: Could not parse {file}: {e}")
                continue
    
    except FileNotFoundError:
        print(f"Error: Manifest path not found: {manifest_path}")
        sys.exit(1)
    
    return manifests


def extract_workload_resources(manifests: List[Dict[str, Any]]) -> List[Dict[str, Any]]:
    """Extract workload resource requirements from manifests"""
    workloads = []
    
    for manifest in manifests:
        kind = manifest.get('kind', '')
        metadata = manifest.get('metadata', {})
        spec = manifest.get('spec', {})
        
        # Extract resource requests/limits
        resources = {}
        
        if kind in ['Deployment', 'StatefulSet', 'DaemonSet', 'Job']:
            template = spec.get('template', {})
            pod_spec = template.get('spec', {})
            containers = pod_spec.get('containers', [])
            
            for container in containers:
                container_name = container.get('name', 'unknown')
                container_resources = container.get('resources', {})
                
                requests = container_resources.get('requests', {})
                limits = container_resources.get('limits', {})
                
                # Use limits if available, otherwise requests
                cpu = limits.get('cpu', requests.get('cpu', '100m'))
                memory = limits.get('memory', requests.get('memory', '128Mi'))
                
 # Convert to standard units
                cpu_cores = parse_cpu(cpu)
                memory_gb = parse_memory(memory)
                
                resources[container_name] = {
                    'cpu_cores': cpu_cores,
                    'memory_gb': memory_gb
                }
            
            # Get replica count
            replicas = spec.get('replicas', 1)
            
            workloads.append({
                'type': kind.lower(),
                'name': metadata.get('name', 'unknown'),
                'namespace': metadata.get('namespace', 'default'),
                'replicas': replicas,
                'resources': resources
            })
        
        elif kind == 'PersistentVolumeClaim':
            storage = spec.get('resources', {}).get('requests', {}).get('storage', '10Gi')
            storage_gb = parse_memory(storage)
            
            workloads.append({
                'type': 'pvc',
                'name': metadata.get('name', 'unknown'),
                'namespace': metadata.get('namespace', 'default'),
                'replicas': 1,
                'resources': {
                    'storage': storage_gb
                }
            })
    
    return workloads


def parse_cpu(cpu_str: str) -> float:
    """Parse CPU string to cores"""
    cpu_str = cpu_str.strip().lower()
    
    if cpu_str.endswith('m'):
        return float(cpu_str[:-1]) / 1000.0
    else:
        return float(cpu_str)


def parse_memory(memory_str: str) -> float:
    """Parse memory string to GB"""
    memory_str = memory_str.strip().upper()
    
    if memory_str.endswith('GI'):
        return float(memory_str[:-2])
    elif memory_str.endswith('G'):
        return float(memory_str[:-1])
    elif memory_str.endswith('MI'):
        return float(memory_str[:-2]) / 1024.0
    elif memory_str.endswith('M'):
        return float(memory_str[:-1]) / 1024.0
    else:
        return float(memory_str) / (1024 * 1024 * 1024)


def predict_workload_costs(workloads: List[Dict[str, Any]]) -> CostPrediction:
    """Predict costs for Kubernetes workloads"""
    workload_costs = []
    total_cpu_cost = 0.0
    total_memory_cost = 0.0
    total_storage_cost = 0.0
    
    for workload in workloads:
        workload_type = workload['type']
        workload_name = workload['name']
        namespace = workload['namespace']
        replicas = workload['replicas']
        resources = workload['resources']
        
        # Calculate resource costs
        cpu_cores = 0.0
        memory_gb = 0.0
        storage_gb = 0.0
        
        if 'storage' in resources:
            storage_gb = resources['storage']
        else:
            for container_name, container_resources in resources.items():
                cpu_cores += container_resources['cpu_cores']
                memory_gb += container_resources['memory_gb']
        
        # Scale by replicas
        cpu_cores *= replicas
        memory_gb *= replicas
        storage_gb *= replicas
        
        # Calculate costs
        cost_breakdown = MockKubernetesPricing.calculate_resource_cost(
            cpu_cores, memory_gb, storage_gb
        )
        
        workload_costs.append(WorkloadCost(
            workload_type=workload_type,
            workload_name=workload_name,
            namespace=namespace,
            monthly_cost=cost_breakdown['total_cost'],
            cpu_cost=cost_breakdown['cpu_cost'],
            memory_cost=cost_breakdown['memory_cost'],
            storage_cost=cost_breakdown['storage_cost']
        ))
        
        total_cpu_cost += cost_breakdown['cpu_cost']
        total_memory_cost += cost_breakdown['memory_cost']
        total_storage_cost += cost_breakdown['storage_cost']
    
    total_monthly = total_cpu_cost + total_memory_cost + total_storage_cost
    
    return CostPrediction(
        total_monthly_cost=total_monthly,
        total_cpu_cost=total_cpu_cost,
        total_memory_cost=total_memory_cost,
        total_storage_cost=total_storage_cost,
        workloads=workload_costs,
        currency="USD"
    )


def format_currency(amount: float) -> str:
    """Format amount as currency"""
    return f"${amount:.2f}"


def generate_report(cost_prediction: CostPrediction, output_format: str = "json") -> str:
    """Generate cost report in specified format"""
    if output_format == "json":
        report = {
            "totalMonthlyCost": cost_prediction.total_monthly_cost,
            "totalCpuCost": cost_prediction.total_cpu_cost,
            "totalMemoryCost": cost_prediction.total_memory_cost,
            "totalStorageCost": cost_prediction.total_storage_cost,
            "currency": cost_prediction.currency,
            "workloads": [
                {
                    "type": w.workload_type,
                    "name": w.workload_name,
                    "namespace": w.namespace,
                    "monthlyCost": w.monthly_cost,
                    "cpuCost": w.cpu_cost,
                    "memoryCost": w.memory_cost,
                    "storageCost": w.storage_cost
                }
                for w in cost_prediction.workloads
            ]
        }
        return json.dumps(report, indent=2)
    
    elif output_format == "text":
        lines = [
            "=== Kubernetes Workload Cost Prediction ===",
            f"Total Monthly Cost: {format_currency(cost_prediction.total_monthly_cost)}",
            f"  CPU Cost: {format_currency(cost_prediction.total_cpu_cost)}",
            f"  Memory Cost: {format_currency(cost_prediction.total_memory_cost)}",
            f"  Storage Cost: {format_currency(cost_prediction.total_storage_cost)}",
            f"Currency: {cost_prediction.currency}",
            "",
            "Workload Breakdown:"
        ]
        
        for workload in cost_prediction.workloads:
            lines.append(f"  - {workload.workload_type}/{workload.workload_name} ({workload.namespace})")
            lines.append(f"    Monthly: {format_currency(workload.monthly_cost)}")
            lines.append(f"    CPU: {format_currency(workload.cpu_cost)}")
            lines.append(f"    Memory: {format_currency(workload.memory_cost)}")
            lines.append(f"    Storage: {format_currency(workload.storage_cost)}")
        
        return "\n".join(lines)
    
    else:
        raise ValueError(f"Unknown output format: {output_format}")


def main():
    parser = argparse.ArgumentParser(description="Mock kubectl-cost for Kubernetes cost prediction")
    parser.add_argument("--path", required=True, help="Path to Kubernetes manifests (file or directory)")
    parser.add_argument("--format", default="json", choices=["json", "text"], help="Output format")
    parser.add_argument("--out-file", help="Output file path")
    
    args = parser.parse_args()
    
    # Parse Kubernetes manifests
    manifests = parse_kubernetes_manifests(args.path)
    
    # Extract workload resources
    workloads = extract_workload_resources(manifests)
    
    # Predict costs
    cost_prediction = predict_workload_costs(workloads)
    
    # Generate report
    report = generate_report(cost_prediction, output_format=args.format)
    
    # Output report
    if args.out_file:
        with open(args.out_file, 'w') as f:
            f.write(report)
        print(f"Cost prediction written to {args.out_file}")
    else:
        print(report)


if __name__ == "__main__":
    main()
