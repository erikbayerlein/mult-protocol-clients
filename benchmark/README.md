# Benchmark Suite

A complete benchmarking system to compare the performance of string, JSON, and Protobuf clients across all operations.

## Installation

### Requirements

- Go 1.24+
- Python 3.7+ (for chart generation)
- pip (to install Python dependencies)

### Setup

1. Install Python dependencies:

```bash
pip install matplotlib numpy
```

2. Navigate to the benchmark directory:

```bash
cd benchmark
```

## Quick Start

### Using Makefile (Recommended)

```bash
# Run benchmarks with default settings (5 iterations)
make benchmark

# Run benchmarks with verbose output
make benchmark-verbose

# Run quick benchmarks (3 iterations)
make benchmark-quick

# Run comprehensive benchmarks (10 iterations)
make benchmark-detailed

# Generate charts
make chart

# Clean results
make clean

# Show help
make help
```

### Using Scripts

```bash
# With 5 iterations per operation (default)
./run_benchmarks.sh

# With 10 iterations per operation
./run_benchmarks.sh -i 10

# With verbose output
./run_benchmarks.sh -v

# With verbose output and 10 iterations
./run_benchmarks.sh -i 10 -v
```

### Using Go

```bash
# With 5 iterations per operation (default)
go run benchmark.go

# With 10 iterations per operation
go run benchmark.go -iterations 10

# With verbose output
go run benchmark.go -verbose

# Specify output directory
go run benchmark.go -output ./my_results
```

### Generate Charts

After running the benchmarks:

```bash
python3 generate_charts.py ./benchmark_results
```

### Output

The following files will be generated in the `benchmark_results/` directory:

- **benchmark_results.csv** - Raw data in CSV format
- **benchmark_results.json** - Raw data in JSON format with metadata
- **report.txt** - Formatted statistics report in text
- **comparison_by_operation.png** - Average time per operation (for each client)
- **client_comparison.png** - Comparison of all clients across all operations
- **distribution.png** - Box plots showing distribution of response times
- **success_rate.png** - Success rate per operation and client

## Operations Tested

The benchmarks test the following operations on each client (string, JSON, and Protobuf):

1. **auth** - User authentication (login)
2. **echo** - Echo text message
3. **sum/soma** - Sum of numbers (list: 1,2,3,4,5,6,7,8,9,10)
4. **timestamp** - Get server timestamp
5. **status** - Check server status
6. **history/historico** - Operation history (limit of 5)
7. **logout** - User logout

## Test Flow

For each client, the flow is:

```
1. AUTH (N iterations)
   ↓
2. ECHO (N iterations) → login (re-enter)
   ↓
3. SUM (N iterations) → login (re-enter)
   ↓
4. TIMESTAMP (N iterations) → login (re-enter)
   ↓
5. STATUS (N iterations) → login (re-enter)
   ↓
6. HISTORY (N iterations) → login (re-enter)
   ↓
7. LOGOUT (N iterations)
```

Where N is the number of iterations specified.

## Project Structure

```
benchmark/
├── benchmark.go           # Main benchmarking code
├── generate_charts.py     # Script to generate charts
├── run_benchmarks.sh      # Shell script to run benchmarks
└── README.md             # This file
```

## Results Interpretation

### Collected Metrics

For each operation and client:

- **Count** - Number of iterations executed
- **Duration (ms)** - Execution time in milliseconds
- **Avg (ms)** - Average execution time
- **Min (ms)** - Minimum time observed
- **Max (ms)** - Maximum time observed
- **Std Dev** - Standard deviation of time distribution
- **Success Rate (%)** - Percentage of successful executions

### Generated Charts

1. **comparison_by_operation.png** - Shows average time for each operation on each client
   - Useful for identifying which operation is slowest on each protocol

2. **client_comparison.png** - Compares all clients side by side
   - Shows which protocol is fastest for each operation

3. **distribution.png** - Box plots show variability and outliers
   - Helps identify consistency and outliers

4. **success_rate.png** - Success rate (should be ~100%)
   - Indicates reliability issues
