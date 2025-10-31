package internal

import (
	"testing"

	"go.uber.org/zap"
)

func BenchmarkSolverLoad(b *testing.B) {
	// Use a real or test logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create a solver pointing to the folder with your test files
	folderPath := "../assets"
	solver := NewSolver(logger, folderPath)

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		if err := solver.Load(); err != nil {
			b.Fatalf("Load failed: %v", err)
		}
	}
}
