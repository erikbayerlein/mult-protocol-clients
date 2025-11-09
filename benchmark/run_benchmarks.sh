#!/bin/bash

# Benchmark runner script
# This script runs the benchmarks and generates charts

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
RESULTS_DIR="${SCRIPT_DIR}/benchmark_results"

echo "═══════════════════════════════════════════════════════════════"
echo "Multi-Protocol Clients Benchmark Suite"
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Parse arguments
ITERATIONS=5
VERBOSE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -i|--iterations)
            ITERATIONS="$2"
            shift 2
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -h|--help)
            echo "Usage: ./run_benchmarks.sh [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  -i, --iterations N    Number of iterations per operation (default: 5)"
            echo "  -v, --verbose         Show detailed output"
            echo "  -h, --help            Show this help message"
            echo ""
            echo "Example:"
            echo "  ./run_benchmarks.sh -i 10 -v"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

echo "Configuration:"
echo "  Iterations: $ITERATIONS"
echo "  Verbose: $VERBOSE"
echo "  Results directory: $RESULTS_DIR"
echo ""

# Run benchmarks
echo "Starting benchmarks..."
cd "$SCRIPT_DIR"

if [ "$VERBOSE" = true ]; then
    go run benchmark.go -iterations=$ITERATIONS -output=$RESULTS_DIR -verbose
else
    go run benchmark.go -iterations=$ITERATIONS -output=$RESULTS_DIR
fi

echo ""
echo "✓ Benchmarks completed"
echo ""

# Generate charts if Python is available
if command -v python3 &> /dev/null; then
    echo "Generating charts..."
    python3 generate_charts.py "$RESULTS_DIR"
    echo ""
    echo "✓ Charts generated successfully"
    
    # Generate detailed analysis
    echo ""
    echo "Running detailed analysis..."
    python3 analyze_results.py "$RESULTS_DIR"
    echo ""
    echo "✓ Analysis completed"
else
    echo "⚠ Python 3 not found. Skipping chart generation."
    echo "  Install Python 3 and matplotlib to generate charts:"
    echo "  pip install matplotlib numpy"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "Results saved to: $RESULTS_DIR"
echo "═══════════════════════════════════════════════════════════════"
