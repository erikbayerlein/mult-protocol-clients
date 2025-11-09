.PHONY: benchmark benchmark-verbose benchmark-quick chart clean help

# Default target
help:
	@echo "╔════════════════════════════════════════════════════════════════╗"
	@echo "║                    Multi-Protocol Clients                      ║"
	@echo "╚════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "Available commands:"
	@echo ""
	@echo "  make run            	  Run program"
	@echo "  make benchmark           Run benchmarks with 5 iterations (default)"
	@echo "  make benchmark-quick     Run benchmarks with 3 iterations (faster)"
	@echo "  make benchmark-verbose   Run benchmarks with verbose output"
	@echo "  make benchmark-detailed  Run benchmarks with 10 iterations (detailed)"
	@echo "  make chart               Generate charts from latest results"
	@echo "  make clean               Remove all benchmark results"
	@echo "  make help                Show this help message"
	@echo ""
	@echo "Example:"
	@echo "  make benchmark-detailed && make chart"
	@echo ""

run:
	@go run .

benchmark:
	@cd benchmark && ./run_benchmarks.sh -i 5

benchmark-quick:
	@cd benchmark && ./run_benchmarks.sh -i 3

benchmark-verbose:
	@cd benchmark && ./run_benchmarks.sh -i 5 -v

benchmark-detailed:
	@cd benchmark && ./run_benchmarks.sh -i 10

chart:
	@if [ -f benchmark/benchmark_results/benchmark_results.csv ]; then \
		python3 benchmark/generate_charts.py benchmark/benchmark_results; \
	else \
		echo "Error: benchmark_results.csv not found. Run 'make benchmark' first."; \
		exit 1; \
	fi

clean:
	@rm -rf benchmark/benchmark_results
	@echo "✓ Benchmark results cleaned"

.DEFAULT_GOAL := help
