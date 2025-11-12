#!/usr/bin/env python3
"""
Benchmark visualization script
Generates charts comparing performance across different clients and operations
"""

import csv
import json
import matplotlib.pyplot as plt
import numpy as np
from pathlib import Path
from collections import defaultdict
import sys

def load_csv_results(csv_file):
    """Load benchmark results from CSV file"""
    results = []
    with open(csv_file, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            results.append({
                'client': row['Client'],
                'operation': row['Operation'],
                'duration_ms': float(row['DurationMs']),
                'mem_alloc_before': int(row['MemAllocBefore']) if 'MemAllocBefore' in row else 0,
                'mem_alloc_after': int(row['MemAllocAfter']) if 'MemAllocAfter' in row else 0,
                'mem_alloc_delta': int(row['MemAllocDelta']) if 'MemAllocDelta' in row else 0,
                'success': row['Success'].lower() == 'true',
                'error': row['Error']
            })
    return results

def aggregate_by_client_operation(results):
    """Aggregate results by client and operation"""
    aggregated = defaultdict(lambda: defaultdict(list))
    
    for result in results:
        client = result['client']
        operation = result['operation']
        aggregated[client][operation].append(result['duration_ms'])
    
    return aggregated

def calculate_stats(durations):
    """Calculate statistics for a list of durations"""
    if not durations:
        return None
    
    return {
        'mean': np.mean(durations),
        'median': np.median(durations),
        'std': np.std(durations),
        'min': np.min(durations),
        'max': np.max(durations),
        'count': len(durations)
    }

def plot_comparison_by_operation(results, output_dir):
    """Create comparison charts for each operation across clients"""
    aggregated = aggregate_by_client_operation(results)
    
    # Get unique operations
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    clients = sorted(aggregated.keys())
    
    # Create a figure for average execution time comparison
    fig, axes = plt.subplots(2, 3, figsize=(15, 10))
    fig.suptitle('Benchmark Results - Average Execution Time by Operation', fontsize=16, fontweight='bold')
    axes = axes.flatten()
    
    for idx, operation in enumerate(operations):
        if idx >= len(axes):
            break
        
        ax = axes[idx]
        client_times = []
        client_labels = []
        
        for client in clients:
            if operation in aggregated[client]:
                durations = aggregated[client][operation]
                stats = calculate_stats(durations)
                if stats:
                    client_times.append(stats['mean'])
                    client_labels.append(client.upper())
        
        if client_times:
            colors = ['#FF6B6B', '#4ECDC4', '#45B7D1']
            bars = ax.bar(client_labels, client_times, color=colors[:len(client_labels)], alpha=0.8, edgecolor='black')
            ax.set_ylabel('Time (ms)', fontweight='bold')
            ax.set_title(f'{operation.upper()}', fontweight='bold')
            ax.grid(axis='y', alpha=0.3)
            
            # Add value labels on bars
            for bar in bars:
                height = bar.get_height()
                ax.text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.2f}ms',
                       ha='center', va='bottom', fontsize=9)
    
    # Hide unused subplots
    for idx in range(len(operations), len(axes)):
        axes[idx].set_visible(False)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/comparison_by_operation.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: comparison_by_operation.png")
    plt.close()

def plot_client_comparison(results, output_dir):
    """Create comparison charts for each client across operations"""
    aggregated = aggregate_by_client_operation(results)
    
    # Get unique operations
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    clients = sorted(aggregated.keys())
    
    # Create grouped bar chart
    fig, ax = plt.subplots(figsize=(14, 8))
    
    x = np.arange(len(operations))
    width = 0.25
    colors = ['#FF6B6B', '#4ECDC4', '#45B7D1']
    
    for idx, client in enumerate(clients):
        means = []
        for operation in operations:
            if operation in aggregated[client]:
                durations = aggregated[client][operation]
                stats = calculate_stats(durations)
                if stats:
                    means.append(stats['mean'])
                else:
                    means.append(0)
            else:
                means.append(0)
        
        ax.bar(x + idx*width, means, width, label=client.upper(), color=colors[idx], alpha=0.8, edgecolor='black')
    
    ax.set_xlabel('Operations', fontweight='bold', fontsize=12)
    ax.set_ylabel('Time (ms)', fontweight='bold', fontsize=12)
    ax.set_title('Performance Comparison: All Operations Across Clients', fontsize=14, fontweight='bold')
    ax.set_xticks(x + width)
    ax.set_xticklabels([op.upper() for op in operations])
    ax.legend(loc='upper left', fontsize=10)
    ax.grid(axis='y', alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/client_comparison.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: client_comparison.png")
    plt.close()

def plot_distribution(results, output_dir):
    """Create box plots showing distribution of execution times"""
    aggregated = aggregate_by_client_operation(results)
    
    # Get unique operations
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    clients = sorted(aggregated.keys())
    
    fig, axes = plt.subplots(1, len(clients), figsize=(5*len(clients), 6), sharey=False)
    if len(clients) == 1:
        axes = [axes]
    
    fig.suptitle('Distribution of Execution Times by Operation', fontsize=14, fontweight='bold')
    
    for idx, client in enumerate(clients):
        ax = axes[idx]
        data = []
        labels = []
        
        for operation in operations:
            if operation in aggregated[client]:
                data.append(aggregated[client][operation])
                labels.append(operation.upper())
        
        if data:
            bp = ax.boxplot(data, labels=labels, patch_artist=True)
            
            for patch in bp['boxes']:
                patch.set_facecolor('#4ECDC4')
                patch.set_alpha(0.7)
            
            ax.set_ylabel('Time (ms)', fontweight='bold')
            ax.set_title(f'{client.upper()} Client', fontweight='bold')
            ax.grid(axis='y', alpha=0.3)
            plt.setp(ax.xaxis.get_majorticklabels(), rotation=45, ha='right')
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/distribution.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: distribution.png")
    plt.close()

def plot_success_rate(results, output_dir):
    """Create chart showing success rates"""
    aggregated = defaultdict(lambda: defaultdict(lambda: {'success': 0, 'total': 0}))
    
    for result in results:
        key = (result['client'], result['operation'])
        aggregated[key[0]][key[1]]['total'] += 1
        if result['success']:
            aggregated[key[0]][key[1]]['success'] += 1
    
    # Get unique operations
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    clients = sorted(aggregated.keys())
    
    fig, ax = plt.subplots(figsize=(12, 6))
    
    x = np.arange(len(operations))
    width = 0.25
    colors = ['#FF6B6B', '#4ECDC4', '#45B7D1']
    
    for idx, client in enumerate(clients):
        rates = []
        for operation in operations:
            if operation in aggregated[client]:
                data = aggregated[client][operation]
                rate = (data['success'] / data['total'] * 100) if data['total'] > 0 else 0
                rates.append(rate)
            else:
                rates.append(0)
        
        ax.bar(x + idx*width, rates, width, label=client.upper(), color=colors[idx], alpha=0.8, edgecolor='black')
    
    ax.set_xlabel('Operations', fontweight='bold', fontsize=12)
    ax.set_ylabel('Success Rate (%)', fontweight='bold', fontsize=12)
    ax.set_title('Success Rate by Operation and Client', fontsize=14, fontweight='bold')
    ax.set_xticks(x + width)
    ax.set_xticklabels([op.upper() for op in operations])
    ax.set_ylim(0, 105)
    ax.legend(loc='lower right', fontsize=10)
    ax.grid(axis='y', alpha=0.3)
    
    # Add percentage labels
    for container in ax.containers:
        ax.bar_label(container, fmt='%.0f%%', fontsize=8)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/success_rate.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: success_rate.png")
    plt.close()

def generate_report(results, output_dir):
    """Generate a text report with statistics"""
    aggregated = aggregate_by_client_operation(results)
    
    report = []
    report.append("=" * 100)
    report.append("BENCHMARK REPORT")
    report.append("=" * 100)
    report.append("")
    
    clients = sorted(aggregated.keys())
    
    for client in clients:
        report.append(f"\n--- {client.upper()} CLIENT ---")
        report.append("-" * 80)
        report.append(f"{'Operation':<20} {'Avg (ms)':<15} {'Min (ms)':<15} {'Max (ms)':<15} {'StdDev':<15}")
        report.append("-" * 80)
        
        for operation in sorted(aggregated[client].keys()):
            durations = aggregated[client][operation]
            stats = calculate_stats(durations)
            if stats:
                report.append(
                    f"{operation:<20} {stats['mean']:<15.4f} {stats['min']:<15.4f} "
                    f"{stats['max']:<15.4f} {stats['std']:<15.4f}"
                )
    
    report.append("\n" + "=" * 100)
    
    report_text = "\n".join(report)
    print(report_text)
    
    with open(f'{output_dir}/report.txt', 'w') as f:
        f.write(report_text)
    
    print(f"✓ Saved: report.txt")

def plot_memory_by_operation(results, output_dir):
    """Create memory comparison charts for each operation across clients"""
    aggregated = aggregate_by_client_operation(results)
    
    # Get unique operations
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    # Extract memory data
    memory_aggregated = defaultdict(lambda: defaultdict(list))
    for result in results:
        client = result['client']
        operation = result['operation']
        memory_aggregated[client][operation].append(result['mem_alloc_delta'])
    
    clients = sorted(memory_aggregated.keys())
    
    # Create a figure for memory comparison
    fig, axes = plt.subplots(2, 3, figsize=(15, 10))
    fig.suptitle('Memory Allocation by Operation (in KB)', fontsize=16, fontweight='bold')
    axes = axes.flatten()
    
    for idx, operation in enumerate(operations):
        if idx >= len(axes):
            break
        
        ax = axes[idx]
        client_memory = []
        client_labels = []
        
        for client in clients:
            if operation in memory_aggregated[client]:
                memory_vals = memory_aggregated[client][operation]
                if memory_vals:
                    avg_memory_kb = np.mean(memory_vals) / 1024.0
                    client_memory.append(avg_memory_kb)
                    client_labels.append(client.upper())
        
        if client_memory:
            colors = ['#FF6B6B', '#4ECDC4', '#45B7D1']
            bars = ax.bar(client_labels, client_memory, color=colors[:len(client_labels)], alpha=0.8, edgecolor='black')
            ax.set_ylabel('Memory (KB)', fontweight='bold')
            ax.set_title(f'{operation.upper()}', fontweight='bold')
            ax.grid(axis='y', alpha=0.3)
            
            # Add value labels on bars
            for bar in bars:
                height = bar.get_height()
                ax.text(bar.get_x() + bar.get_width()/2., height,
                       f'{height:.2f}KB',
                       ha='center', va='bottom', fontsize=9)
    
    # Hide unused subplots
    for idx in range(len(operations), len(axes)):
        axes[idx].set_visible(False)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/memory_by_operation.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: memory_by_operation.png")
    plt.close()

def plot_memory_client_comparison(results, output_dir):
    """Create memory comparison charts for each client across operations"""
    # Extract memory data
    memory_aggregated = defaultdict(lambda: defaultdict(list))
    for result in results:
        client = result['client']
        operation = result['operation']
        memory_aggregated[client][operation].append(result['mem_alloc_delta'])
    
    # Get unique operations
    operations = set()
    for client_data in memory_aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    clients = sorted(memory_aggregated.keys())
    
    # Create grouped bar chart for memory
    fig, ax = plt.subplots(figsize=(14, 8))
    
    x = np.arange(len(operations))
    width = 0.25
    colors = ['#FF6B6B', '#4ECDC4', '#45B7D1']
    
    for idx, client in enumerate(clients):
        memory_means = []
        for operation in operations:
            if operation in memory_aggregated[client]:
                memory_vals = memory_aggregated[client][operation]
                if memory_vals:
                    avg_memory_kb = np.mean(memory_vals) / 1024.0
                    memory_means.append(avg_memory_kb)
                else:
                    memory_means.append(0)
            else:
                memory_means.append(0)
        
        ax.bar(x + idx*width, memory_means, width, label=client.upper(), color=colors[idx], alpha=0.8, edgecolor='black')
    
    ax.set_xlabel('Operations', fontweight='bold', fontsize=12)
    ax.set_ylabel('Memory Allocation (KB)', fontweight='bold', fontsize=12)
    ax.set_title('Memory Allocation Comparison: All Operations Across Clients', fontsize=14, fontweight='bold')
    ax.set_xticks(x + width)
    ax.set_xticklabels([op.upper() for op in operations])
    ax.legend(loc='upper left', fontsize=10)
    ax.grid(axis='y', alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/memory_client_comparison.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: memory_client_comparison.png")
    plt.close()

def plot_time_vs_memory(results, output_dir):
    """Create scatter plot comparing execution time vs memory usage"""
    fig, ax = plt.subplots(figsize=(12, 8))
    
    colors_map = {'string': '#FF6B6B', 'json': '#4ECDC4', 'proto': '#45B7D1'}
    
    for client in sorted(set(r['client'] for r in results)):
        client_results = [r for r in results if r['client'] == client]
        
        times = [r['duration_ms'] for r in client_results]
        memories = [r['mem_alloc_delta'] / 1024.0 for r in client_results]
        operations = [r['operation'] for r in client_results]
        
        ax.scatter(times, memories, label=client.upper(), s=100, alpha=0.6, 
                  color=colors_map.get(client, '#999999'), edgecolor='black', linewidth=1.5)
    
    ax.set_xlabel('Execution Time (ms)', fontweight='bold', fontsize=12)
    ax.set_ylabel('Memory Allocation (KB)', fontweight='bold', fontsize=12)
    ax.set_title('Execution Time vs Memory Allocation', fontsize=14, fontweight='bold')
    ax.legend(loc='upper left', fontsize=10)
    ax.grid(True, alpha=0.3)
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/time_vs_memory.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: time_vs_memory.png")
    plt.close()

def plot_memory_distribution(results, output_dir):
    """Create box plots showing distribution of memory allocations"""
    # Extract memory data
    memory_aggregated = defaultdict(lambda: defaultdict(list))
    for result in results:
        client = result['client']
        operation = result['operation']
        memory_aggregated[client][operation].append(result['mem_alloc_delta'] / 1024.0)
    
    # Get unique operations
    operations = set()
    for client_data in memory_aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    clients = sorted(memory_aggregated.keys())
    
    fig, axes = plt.subplots(1, len(clients), figsize=(5*len(clients), 6), sharey=False)
    if len(clients) == 1:
        axes = [axes]
    
    fig.suptitle('Distribution of Memory Allocations by Operation', fontsize=14, fontweight='bold')
    
    for idx, client in enumerate(clients):
        ax = axes[idx]
        data = []
        labels = []
        
        for operation in operations:
            if operation in memory_aggregated[client]:
                data.append(memory_aggregated[client][operation])
                labels.append(operation.upper())
        
        if data:
            bp = ax.boxplot(data, labels=labels, patch_artist=True)
            
            for patch in bp['boxes']:
                patch.set_facecolor('#4ECDC4')
                patch.set_alpha(0.7)
            
            ax.set_ylabel('Memory (KB)', fontweight='bold')
            ax.set_title(f'{client.upper()} Client', fontweight='bold')
            ax.grid(axis='y', alpha=0.3)
            plt.setp(ax.xaxis.get_majorticklabels(), rotation=45, ha='right')
    
    plt.tight_layout()
    plt.savefig(f'{output_dir}/memory_distribution.png', dpi=300, bbox_inches='tight')
    print(f"✓ Saved: memory_distribution.png")
    plt.close()

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 generate_charts.py <benchmark_results_dir>")
        sys.exit(1)
    
    results_dir = Path(sys.argv[1])
    csv_file = results_dir / 'benchmark_results.csv'
    
    if not csv_file.exists():
        print(f"Error: CSV file not found at {csv_file}")
        sys.exit(1)
    
    print(f"Loading results from {csv_file}...")
    results = load_csv_results(csv_file)
    print(f"✓ Loaded {len(results)} results\n")
    
    print("Generating charts...")
    plot_comparison_by_operation(results, results_dir)
    plot_client_comparison(results, results_dir)
    plot_distribution(results, results_dir)
    plot_success_rate(results, results_dir)
    plot_memory_by_operation(results, results_dir)
    plot_memory_client_comparison(results, results_dir)
    plot_time_vs_memory(results, results_dir)
    plot_memory_distribution(results, results_dir)
    generate_report(results, results_dir)
    
    print(f"\n✓ All charts saved to {results_dir}")

if __name__ == '__main__':
    main()
