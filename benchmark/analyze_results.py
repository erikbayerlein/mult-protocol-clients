#!/usr/bin/env python3
"""
Benchmark Data Analyzer
Provides detailed analysis and comparison of benchmark results
"""

import csv
import json
from pathlib import Path
from collections import defaultdict
import statistics
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
                'success': row['Success'].lower() == 'true',
                'error': row['Error']
            })
    return results

def aggregate_by_client_operation(results):
    """Aggregate results by client and operation"""
    aggregated = defaultdict(lambda: defaultdict(list))
    
    for result in results:
        aggregated[result['client']][result['operation']].append(result['duration_ms'])
    
    return aggregated

def calculate_stats(durations):
    """Calculate comprehensive statistics for a list of durations"""
    if not durations:
        return None
    
    sorted_dur = sorted(durations)
    
    return {
        'count': len(durations),
        'mean': statistics.mean(durations),
        'median': statistics.median(durations),
        'stdev': statistics.stdev(durations) if len(durations) > 1 else 0,
        'min': min(durations),
        'max': max(durations),
        'q1': sorted_dur[len(sorted_dur)//4],
        'q3': sorted_dur[3*len(sorted_dur)//4],
        'iqr': sorted_dur[3*len(sorted_dur)//4] - sorted_dur[len(sorted_dur)//4],
    }

def print_detailed_report(results):
    """Print a detailed analysis report"""
    aggregated = aggregate_by_client_operation(results)
    
    print("\n" + "╔" + "═" * 118 + "╗")
    print("║" + " " * 45 + "DETAILED BENCHMARK ANALYSIS" + " " * 47 + "║")
    print("╚" + "═" * 118 + "╝\n")
    
    clients = sorted(aggregated.keys())
    
    for client in clients:
        print(f"\n┌─ {client.upper()} CLIENT " + "─" * 100)
        
        operations = sorted(aggregated[client].keys())
        
        for operation in operations:
            durations = aggregated[client][operation]
            stats = calculate_stats(durations)
            
            if stats:
                print(f"\n  {operation.upper()}")
                print(f"  ├─ Iterations:  {stats['count']}")
                print(f"  ├─ Mean:        {stats['mean']:.4f} ms")
                print(f"  ├─ Median:      {stats['median']:.4f} ms")
                print(f"  ├─ Std Dev:     {stats['stdev']:.4f} ms")
                print(f"  ├─ Min:         {stats['min']:.4f} ms")
                print(f"  ├─ Max:         {stats['max']:.4f} ms")
                print(f"  ├─ Q1 (25%):    {stats['q1']:.4f} ms")
                print(f"  ├─ Q3 (75%):    {stats['q3']:.4f} ms")
                print(f"  └─ IQR:         {stats['iqr']:.4f} ms")

def print_comparison_table(results):
    """Print comparison table for all operations"""
    aggregated = aggregate_by_client_operation(results)
    clients = sorted(aggregated.keys())
    
    # Get all operations
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    print("\n" + "┌" + "─" * 130 + "┐")
    print("│" + " " * 50 + "MEAN EXECUTION TIME COMPARISON (ms)" + " " * 45 + "│")
    print("├" + "─" * 130 + "┤")
    
    # Header
    header = "│ {:<20} │".format("Operation")
    for client in clients:
        header += " {:>15} │".format(client.upper())
    print(header)
    
    print("├" + "─" * 130 + "┤")
    
    # Data rows
    for operation in operations:
        row = "│ {:<20} │".format(operation.upper())
        for client in clients:
            if operation in aggregated[client]:
                durations = aggregated[client][operation]
                mean_time = sum(durations) / len(durations)
                row += " {:>15.4f} │".format(mean_time)
            else:
                row += " {:>15} │".format("N/A")
        print(row)
    
    print("└" + "─" * 130 + "┘")

def print_performance_ranking(results):
    """Print ranking of clients by performance (fastest to slowest)"""
    aggregated = aggregate_by_client_operation(results)
    
    print("\n" + "┌" + "─" * 118 + "┐")
    print("│" + " " * 45 + "PERFORMANCE RANKING BY OPERATION" + " " * 40 + "│")
    print("└" + "─" * 118 + "┘")
    
    operations = set()
    for client_data in aggregated.values():
        operations.update(client_data.keys())
    operations = sorted(operations)
    
    for operation in operations:
        clients_times = []
        for client in sorted(aggregated.keys()):
            if operation in aggregated[client]:
                durations = aggregated[client][operation]
                mean_time = sum(durations) / len(durations)
                clients_times.append((client, mean_time))
        
        clients_times.sort(key=lambda x: x[1])
        
        print(f"\n  {operation.upper()}")
        for rank, (client, time) in enumerate(clients_times, 1):
            medal = "🥇" if rank == 1 else "🥈" if rank == 2 else "🥉" if rank == 3 else f"#{rank}"
            print(f"    {medal}  {client.upper():<10} {time:>10.4f} ms")

def print_summary_statistics(results):
    """Print summary statistics"""
    aggregated = aggregate_by_client_operation(results)
    
    print("\n" + "┌" + "─" * 118 + "┐")
    print("│" + " " * 50 + "SUMMARY STATISTICS" + " " * 50 + "│")
    print("└" + "─" * 118 + "┘")
    
    clients = sorted(aggregated.keys())
    
    for client in clients:
        all_durations = []
        for durations in aggregated[client].values():
            all_durations.extend(durations)
        
        stats = calculate_stats(all_durations)
        
        print(f"\n  {client.upper()}")
        print(f"    Total Operations:  {stats['count']}")
        print(f"    Overall Mean:      {stats['mean']:.4f} ms")
        print(f"    Overall Std Dev:   {stats['stdev']:.4f} ms")
        print(f"    Overall Min:       {stats['min']:.4f} ms")
        print(f"    Overall Max:       {stats['max']:.4f} ms")

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 analyze_results.py <benchmark_results_dir>")
        sys.exit(1)
    
    results_dir = Path(sys.argv[1])
    csv_file = results_dir / 'benchmark_results.csv'
    
    if not csv_file.exists():
        print(f"Error: CSV file not found at {csv_file}")
        sys.exit(1)
    
    print(f"Loading results from {csv_file}...")
    results = load_csv_results(csv_file)
    print(f"✓ Loaded {len(results)} results\n")
    
    # Print analysis
    print_comparison_table(results)
    print_performance_ranking(results)
    print_summary_statistics(results)
    print_detailed_report(results)
    
    print("\n" + "═" * 120 + "\n")

if __name__ == '__main__':
    main()
