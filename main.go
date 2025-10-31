package main

import (
	"antsolver/internal"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	start := time.Now()
	solver := internal.NewSolver(logger, "./assets")
	if err := solver.Load(); err != nil {
		logger.With(zap.Error(err)).Fatal("Fail to load the dictionary")
		return
	}
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Took %.2f seconds to load the dictionary", elapsed.Seconds()))

	server := internal.NewHTTPServer(logger, solver)
	server.Start()
}
