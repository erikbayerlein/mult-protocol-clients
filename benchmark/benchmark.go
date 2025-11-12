package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/erikbayerlein/mult-protocol-clients/internal/auth"
	jc "github.com/erikbayerlein/mult-protocol-clients/json"
	pb "github.com/erikbayerlein/mult-protocol-clients/proto"
	sc "github.com/erikbayerlein/mult-protocol-clients/strings"
)

const (
	host           = "3.88.99.255"
	string_port    = 8080
	json_port      = 8081
	protobuff_port = 8082
	studentID      = 537606
)

type BenchmarkResult struct {
	Client         string
	Operation      string
	Iteration      int
	Duration       time.Duration
	DurationMs     float64
	MemAllocBefore uint64
	MemAllocAfter  uint64
	MemAllocDelta  uint64
	Success        bool
	Error          string
}

type BenchmarkStats struct {
	Client      string
	Operation   string
	Count       int
	TotalTimeMs float64
	AvgTimeMs   float64
	MinTimeMs   float64
	MaxTimeMs   float64
	StdDevMs    float64
	MemAllocAvg uint64
	SuccessRate float64
}

var (
	iterations = flag.Int("iterations", 5, "Number of iterations per operation")
	outputDir  = flag.String("output", "./benchmark_results", "Output directory for results")
	verbose    = flag.Bool("verbose", false, "Verbose output")
)

func main() {
	flag.Parse()

	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting benchmarks...")
	fmt.Printf("Iterations: %d\n", *iterations)
	fmt.Printf("Output directory: %s\n\n", *outputDir)

	// Initialize clients
	stringClient := sc.StringClient{Host: host, Port: string_port}
	jsonClient := jc.JsonClient{Host: host, Port: json_port}
	protoClient := pb.ProtobufClient{Host: host, Port: protobuff_port}

	var results []BenchmarkResult

	// Test each client
	results = append(results, benchmarkClient("string", &stringClient)...)
	results = append(results, benchmarkClient("json", &jsonClient)...)
	results = append(results, benchmarkClient("proto", &protoClient)...)

	// Save results
	if err := saveResultsToCSV(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving CSV: %v\n", err)
		os.Exit(1)
	}

	if err := saveResultsToJSON(results); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving JSON: %v\n", err)
		os.Exit(1)
	}

	// Calculate and print statistics
	stats := calculateStats(results)
	printStats(stats)

	fmt.Printf("\nResults saved to %s\n", *outputDir)
	fmt.Println("✓ Benchmarks completed successfully")
}

func benchmarkClient(clientName string, client interface{}) []BenchmarkResult {
	var results []BenchmarkResult

	fmt.Printf("=== Benchmarking %s client ===\n", strings.ToUpper(clientName))

	// Get token
	token := ""
	authResults := benchmarkAuth(clientName, client)
	results = append(results, authResults...)

	rec, err := auth.LoadToken()
	if err != nil {
		fmt.Printf("Error loading token: %v\n", err)
		return results
	}
	token = rec.Token

	operations := []struct {
		name string
		args []string
	}{
		{"echo", []string{"benchmark test message"}},
		{"sum", []string{"1,2,3,4,5,6,7,8,9,10"}},
		{"timestamp", []string{}},
		{"status", []string{}},
		{"history", []string{"5"}},
	}

	for _, op := range operations {
		fmt.Printf("\n  Testing %s...\n", op.name)

		for i := 0; i < *iterations; i++ {
			result := BenchmarkResult{
				Client:    clientName,
				Operation: op.name,
				Iteration: i + 1,
			}

			runtime.GC()
			memBefore := getMemoryStats()
			result.MemAllocBefore = memBefore.Alloc

			start := time.Now()
			var err error

			switch c := client.(type) {
			case *sc.StringClient:
				_, err = c.DoOperation(op.name, token, parseStringParams(op.name, op.args))
			case *jc.JsonClient:
				err = c.Run(op.name, op.args)
			case *pb.ProtobufClient:
				err = c.Run(op.name, op.args)
			}

			duration := time.Since(start)
			memAfter := getMemoryStats()
			result.MemAllocAfter = memAfter.Alloc
			result.Duration = duration
			result.DurationMs = duration.Seconds() * 1000

			if result.MemAllocAfter >= result.MemAllocBefore {
				result.MemAllocDelta = result.MemAllocAfter - result.MemAllocBefore
			}

			if err != nil {
				result.Success = false
				result.Error = err.Error()
				if *verbose {
					fmt.Printf("    [%d] %s - ERROR: %v (%v)\n", i+1, op.name, err, duration)
				}
			} else {
				result.Success = true
				if *verbose {
					fmt.Printf("    [%d] %s - OK (%v)\n", i+1, op.name, duration)
				}
			}

			results = append(results, result)
		}
	}

	// Logout benchmark
	fmt.Printf("\n  Testing logout...\n")
	for i := 0; i < *iterations; i++ {
		result := BenchmarkResult{
			Client:    clientName,
			Operation: "logout",
			Iteration: i + 1,
		}

		runtime.GC()
		memBefore := getMemoryStats()
		result.MemAllocBefore = memBefore.Alloc

		start := time.Now()
		var err error

		switch c := client.(type) {
		case *sc.StringClient:
			err = c.Logout(token)
		case *jc.JsonClient:
			err = c.Logout(token)
		case *pb.ProtobufClient:
			err = c.Logout(token)
		}

		duration := time.Since(start)
		memAfter := getMemoryStats()
		result.MemAllocAfter = memAfter.Alloc
		result.Duration = duration
		result.DurationMs = duration.Seconds() * 1000

		if result.MemAllocAfter >= result.MemAllocBefore {
			result.MemAllocDelta = result.MemAllocAfter - result.MemAllocBefore
		}

		if err != nil {
			result.Success = false
			result.Error = err.Error()
			if *verbose {
				fmt.Printf("    [%d] logout - ERROR: %v (%v)\n", i+1, err, duration)
			}
		} else {
			result.Success = true
			if *verbose {
				fmt.Printf("    [%d] logout - OK (%v)\n", i+1, duration)
			}
		}

		results = append(results, result)

		// After first logout, login again for next iterations
		if i < *iterations-1 {
			if err := benchmarkLogin(client); err != nil {
				fmt.Printf("Error re-logging in: %v\n", err)
				break
			}
			rec, _ := auth.LoadToken()
			token = rec.Token
		}
	}

	fmt.Printf("\n✓ %s client benchmarks completed\n", strings.ToUpper(clientName))
	return results
}

func benchmarkAuth(clientName string, client interface{}) []BenchmarkResult {
	var results []BenchmarkResult

	fmt.Printf("\n  Testing auth...\n")

	for i := 0; i < *iterations; i++ {
		result := BenchmarkResult{
			Client:    clientName,
			Operation: "auth",
			Iteration: i + 1,
		}

		runtime.GC()
		memBefore := getMemoryStats()
		result.MemAllocBefore = memBefore.Alloc

		start := time.Now()
		var err error

		switch c := client.(type) {
		case *sc.StringClient:
			err = c.Login(studentID)
		case *jc.JsonClient:
			err = c.Login(studentID)
		case *pb.ProtobufClient:
			err = c.Login(studentID)
		}

		duration := time.Since(start)
		memAfter := getMemoryStats()
		result.MemAllocAfter = memAfter.Alloc
		result.Duration = duration
		result.DurationMs = duration.Seconds() * 1000

		if result.MemAllocAfter >= result.MemAllocBefore {
			result.MemAllocDelta = result.MemAllocAfter - result.MemAllocBefore
		}

		if err != nil {
			result.Success = false
			result.Error = err.Error()
			if *verbose {
				fmt.Printf("    [%d] auth - ERROR: %v (%v)\n", i+1, err, duration)
			}
		} else {
			result.Success = true
			if *verbose {
				fmt.Printf("    [%d] auth - OK (%v)\n", i+1, duration)
			}
		}

		results = append(results, result)
	}

	return results
}

func benchmarkLogin(client interface{}) error {
	switch c := client.(type) {
	case *sc.StringClient:
		return c.Login(studentID)
	case *jc.JsonClient:
		return c.Login(studentID)
	case *pb.ProtobufClient:
		return c.Login(studentID)
	}
	return fmt.Errorf("unknown client type")
}

func parseStringParams(operation string, args []string) map[string]interface{} {
	params := make(map[string]interface{})

	switch operation {
	case "echo":
		if len(args) > 0 {
			params["mensagem"] = strings.Join(args, " ")
		}
	case "soma":
		if len(args) > 0 {
			parts := strings.Split(args[0], ",")
			ints := make([]int, 0, len(parts))
			for _, p := range parts {
				n, _ := strconv.Atoi(strings.TrimSpace(p))
				ints = append(ints, n)
			}
			params["nums"] = ints
		}
	case "historico":
		if len(args) > 0 {
			if v, err := strconv.Atoi(args[0]); err == nil {
				params["limite"] = v
			}
		}
	case "status":
		params["detalhado"] = true
	}

	return params
}

func parseParams(operation string, args []string) map[string]interface{} {
	params := make(map[string]interface{})

	switch operation {
	case "echo":
		if len(args) > 0 {
			params["mensagem"] = strings.Join(args, " ")
		}
	case "sum":
		if len(args) > 0 {
			params["nums"] = args[0]
		}
	case "history":
		if len(args) > 0 {
			params["limite"] = args[0]
		}
	case "status":
		params["detalhado"] = "true"
	}

	return params
}

func saveResultsToCSV(results []BenchmarkResult) error {
	filename := fmt.Sprintf("%s/benchmark_results.csv", *outputDir)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Client", "Operation", "Iteration", "DurationMs", "MemAllocBefore", "MemAllocAfter", "MemAllocDelta", "Success", "Error"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, result := range results {
		record := []string{
			result.Client,
			result.Operation,
			strconv.Itoa(result.Iteration),
			fmt.Sprintf("%.4f", result.DurationMs),
			strconv.FormatUint(result.MemAllocBefore, 10),
			strconv.FormatUint(result.MemAllocAfter, 10),
			strconv.FormatUint(result.MemAllocDelta, 10),
			strconv.FormatBool(result.Success),
			result.Error,
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	fmt.Printf("✓ CSV results saved to %s\n", filename)
	return nil
}

func saveResultsToJSON(results []BenchmarkResult) error {
	filename := fmt.Sprintf("%s/benchmark_results.json", *outputDir)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]interface{}{
		"timestamp":  time.Now().Format(time.RFC3339),
		"iterations": *iterations,
		"results":    results,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	fmt.Printf("✓ JSON results saved to %s\n", filename)
	return nil
}

func calculateStats(results []BenchmarkResult) []BenchmarkStats {
	statsMap := make(map[string][]float64)
	memAllocMap := make(map[string][]uint64)
	successMap := make(map[string]int)
	countMap := make(map[string]int)

	// Aggregate data
	for _, result := range results {
		key := result.Client + "_" + result.Operation
		statsMap[key] = append(statsMap[key], result.DurationMs)
		memAllocMap[key] = append(memAllocMap[key], result.MemAllocDelta)
		countMap[key]++
		if result.Success {
			successMap[key]++
		}
	}

	// Calculate statistics
	var stats []BenchmarkStats
	for key, durations := range statsMap {
		parts := strings.Split(key, "_")
		client := parts[0]
		operation := parts[1]

		if len(durations) == 0 {
			continue
		}

		// Calculate metrics
		total := 0.0
		minVal := durations[0]
		maxVal := durations[0]

		for _, d := range durations {
			total += d
			if d < minVal {
				minVal = d
			}
			if d > maxVal {
				maxVal = d
			}
		}

		avgVal := total / float64(len(durations))
		stdDev := calculateStdDev(durations)

		// Calculate memory stats
		memAllocAvg := uint64(0)
		if memAllocVals, ok := memAllocMap[key]; ok && len(memAllocVals) > 0 {
			totalMem := uint64(0)
			for _, m := range memAllocVals {
				totalMem += m
			}
			memAllocAvg = totalMem / uint64(len(memAllocVals))
		}

		successRate := float64(successMap[key]) / float64(countMap[key]) * 100

		stats = append(stats, BenchmarkStats{
			Client:      client,
			Operation:   operation,
			Count:       countMap[key],
			TotalTimeMs: total,
			AvgTimeMs:   avgVal,
			MinTimeMs:   minVal,
			MaxTimeMs:   maxVal,
			StdDevMs:    stdDev,
			MemAllocAvg: memAllocAvg,
			SuccessRate: successRate,
		})
	}

	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Client != stats[j].Client {
			return stats[i].Client < stats[j].Client
		}
		return stats[i].Operation < stats[j].Operation
	})

	return stats
}

func printStats(stats []BenchmarkStats) {
	fmt.Println("\n" + strings.Repeat("=", 120))
	fmt.Println("BENCHMARK STATISTICS")
	fmt.Println(strings.Repeat("=", 120))

	currentClient := ""
	for _, stat := range stats {
		if stat.Client != currentClient {
			currentClient = stat.Client
			fmt.Printf("\n--- %s CLIENT ---\n", strings.ToUpper(stat.Client))
			fmt.Printf("%-15s | %-10s | %-10s | %-10s | %-10s | %-12s | %-10s\n",
				"Operation", "Avg (ms)", "Min (ms)", "Max (ms)", "Success", "MemAlloc(B)", "Count")
			fmt.Println(strings.Repeat("-", 100))
		}

		memAllocKB := float64(stat.MemAllocAvg) / 1024.0
		fmt.Printf("%-15s | %10.4f | %10.4f | %10.4f | %9.1f%% | %12.2f KB | %5d\n",
			stat.Operation,
			stat.AvgTimeMs,
			stat.MinTimeMs,
			stat.MaxTimeMs,
			stat.SuccessRate,
			memAllocKB,
			stat.Count)
	}

	fmt.Println("\n" + strings.Repeat("=", 120))
}

func getMemoryStats() runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m
}

func calculateStdDev(values []float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	mean := 0.0
	for _, v := range values {
		mean += v
	}
	mean /= float64(len(values))

	sumSqDiff := 0.0
	for _, v := range values {
		sumSqDiff += (v - mean) * (v - mean)
	}

	variance := sumSqDiff / float64(len(values)-1)
	return math.Sqrt(variance)
}
