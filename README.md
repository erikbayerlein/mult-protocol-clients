# Multi-Protocol Clients

## Features

🚀 **Three Protocol Implementations**
- String Protocol: Simple text-based communication
- JSON Protocol: Structured data exchange using JSON
- Protobuf Protocol: High-performance binary serialization

⚡ **Interactive CLI**
- Real-time command execution
- Multi-client support
- Session management with token-based authentication

📊 **Benchmark Suite**
- Comprehensive performance testing across all operations
- Automated chart generation

📈 **Seven Core Operations**
- `auth` - Authenticate and obtain session token
- `echo` - Echo text messages from the server
- `sum` - Calculate sum of multiple numbers
- `timestamp` - Get server timestamp information
- `status` - Check server status and details
- `history` - Retrieve operation history
- `logout` - Terminate session and clear token

## Project Structure

```
multi-protocol-clients/
├── main.go                          # Interactive CLI entry point
├── go.mod                           # Go module definition
├── Makefile                         # Build automation
├── benchmark/                       # Benchmark suite
│   ├── benchmark.go                 # Core benchmark program
│   ├── generate_charts.py           # Chart generation from results
│   ├── analyze_results.py           # Statistical analysis
│   ├── run_benchmarks.sh            # Automation script
│   ├── Makefile                     # Benchmark build targets
│   └── README.md                    # Benchmark documentation
├── strings/                         # String protocol client
│   └── client.go                    # String protocol implementation
├── json/                            # JSON protocol client
│   ├── client.go                    # JSON protocol implementation
│   └── dtos.go                      # JSON data structures
├── proto/                           # Protobuf protocol client
│   └── client.go                    # Protobuf protocol implementation
└── internal/                        # Internal packages
    ├── auth/                        # Authentication management
    │   └── auth.go                  # Token storage & retrieval
    ├── tcp/                         # TCP transport layer
    │   └── connection.go            # Connection management
    └── pb/                          # Protocol Buffer definitions
        └── client.pb.go           # Generated protobuf
```

## Architecture

### Protocol Layers

```
┌─────────────────────────────────────────────────────────┐
│                    User Interface                        │
│              (Interactive TUI / CLI)                     │
└──────┬──────────────────────────────────────────────────┘
       │
       ├─────────────────┬────────────────────┐
       │                 │                    │                    
       ▼                 ▼                    ▼                    
┌─────────────┐  ┌──────────────┐  ┌──────────────────┐
│   String    │  │     JSON     │  │    Protobuf      │
│   Client    │  │    Client    │  │     Client       │
└──────┬──────┘  └──────┬───────┘  └─────────┬────────┘
       │                │                    │
       └────────────────┼────────────────────┘
                        │
                        ▼
                ┌──────────────────┐
                │  TCP Transport   │
                │  (Port-specific) │
                └────────┬─────────┘
                         │
                         ▼
                ┌──────────────────┐
                │  Server (Remote) │
                └──────────────────┘
```

### Client Architecture

Each protocol client follows the same interface pattern:

```
┌──────────────────────────────────────┐
│         Protocol Client              │
│  (String/JSON/Protobuf)              │
├──────────────────────────────────────┤
│ • Auth(studentID) -> token           │
│ • Echo(message) -> echo              │
│ • Sum(numbers) -> result             │
│ • Timestamp() -> time_info           │
│ • Status(detailed) -> status_info    │
│ • History(limit) -> operations       │
│ • Logout(token) -> void              │
└──────────┬───────────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│      Internal Auth Manager           │
│  (Token Storage & Validation)        │
└──────────┬───────────────────────────┘
           │
           ▼
┌──────────────────────────────────────┐
│       TCP Connection Layer           │
│  (Network Communication)             │
└──────────────────────────────────────┘
```

## Installation

### Prerequisites

- **Go** 1.24 or higher
- **Python** 3.7+ (for benchmark visualizations)
- **pip** (for Python dependency management)
- **Network access** to server at `3.88.99.255` (ports 8080-8082)

### Setup Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/erikbayerlein/mult-protocol-clients.git
   cd mult-protocol-clients
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Install Python dependencies** (for benchmark suite)
   ```bash
   pip install matplotlib numpy pandas
   ```

4. **Build the project**
   ```bash
   go build -o multi-protocol-clients
   ```

## Usage

### Interactive CLI

```bash
./multi-protocol-clients
```

#### Available Commands

```
login <client> <student_id>        Authenticate with a specific client
whoami                             Show current user and active client
logout                             Logout and clear session
string <operation> [args...]       Run operation with string client
json <operation> [args...]         Run operation with json client
proto <operation> [args...]        Run operation with protobuf client
exit / quit                        Exit the program
clear                              Clear terminal screen
help                               Show command help
```

#### Operation Examples

**Login and Echo**
```
> login string 12345
✓ Successfully logged in as user 12345
Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

> string echo "Hello, World!"
Echo: Hello, World!
Time: 45ms
```

**Sum Numbers**
```
> json sum 10,20,30,40
Result: 100
Time: 32ms
```

**Check Status**
```
> proto status
Status: Active
Clients Connected: 15
Uptime: 48h 32m
```

**History**
```
> string history 5
1. echo @ 14:23:45 - "Hello, World!"
2. sum @ 14:24:12 - [10, 20, 30, 40]
3. timestamp @ 14:25:03 - Success
4. status @ 14:26:18 - Active
5. echo @ 14:27:01 - "Test"
```

### Benchmark Suite

Comprehensive performance testing and comparison of all three clients:

#### Quick Start

```bash
# Run CLI
make run

# Run with default settings (5 iterations per operation)
make benchmark

# Quick benchmark (3 iterations, faster)
make benchmark-quick

# Detailed benchmark (10 iterations, comprehensive)
make benchmark-detailed

# Verbose output for debugging
make benchmark-verbose
```

#### Manual Execution

```bash
cd benchmark
./run_benchmarks.sh -i 5 -v
```

#### Benchmark Options

```
Usage: ./run_benchmarks.sh [OPTIONS]

Options:
  -i, --iterations N    Number of iterations per operation (default: 5)
  -v, --verbose         Enable verbose output during benchmarking
  -h, --help            Show help message

Examples:
  ./run_benchmarks.sh -i 10          # Run with 10 iterations
  ./run_benchmarks.sh -i 5 -v        # Run with verbose output
  ./run_benchmarks.sh                # Run with defaults
```

#### Benchmark Output

The benchmark suite generates the following outputs in the `benchmark/results/` directory:

**CSV Results**
- `benchmark_results.csv` - Raw timing data for all operations

**JSON Results**
- `benchmark_results.json` - Structured results with statistics

**Visualizations** (PNG Charts)
- `comparison_by_operation.png` - Performance by operation type
- `client_comparison.png` - Overall client performance
- `distribution.png` - Response time distribution (box plots)
- `success_rate.png` - Success rates per operation

**Analysis Report**
- Console output with:
  - Summary statistics per client
  - Performance rankings (🥇 🥈 🥉)
  - Operation-specific insights

## Benchmark Results Analysis

The benchmark suite provides detailed performance analysis:

### Performance Metrics

Each operation is measured across:
- **Response Time** (milliseconds)
- **Min/Max** latencies
- **Average** response time
- **Standard Deviation** (consistency)
- **Success Rate** (percentage)

## Configuration

### Server Endpoints

The default server configuration is hardcoded in `main.go`:

```go
const (
    host           = "3.88.99.255"
    string_port    = 8080    // String protocol
    json_port      = 8081    // JSON protocol
    protobuff_port = 8082    // Protobuf protocol
)
```

To use different servers, modify these values in `main.go` and rebuild.

### Typical Use Cases

| Protocol | Best For |
|----------|----------|
| **String** | Development, debugging, simple systems |
| **JSON** | Web services, APIs, general purpose |
| **Protobuf** | High-performance systems, microservices |

## Author

**Erik Bayerlein**

- GitHub: [@erikbayerlein](https://github.com/erikbayerlein)
- Project: Multi-Protocol Clients Benchmark System

---

**Last Updated**: 2025