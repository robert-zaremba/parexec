// Package parexec provides functions/workers to execute calls concurrently
// with options to limit a number of workers/executors.
//
// Check *_test.go files for examples and benchmarks.
// Benchmark results:
//    BenchmarkSemaphore-4               2000000000  0.33 ns/op
//    BenchmarkSimpleRun_with_results-4  2000000000  0.37 ns/op
//    BenchmarkSimpleRun-4               2000000000  0.26 ns/op
package parexec
