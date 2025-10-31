package main

import (
	"antsolver/internal"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const PORT = 3000

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	start := time.Now()
	solver := internal.NewSolver(logger, "../assets")
	solver.Load()
	elapsed := time.Since(start)
	logger.Info(fmt.Sprintf("Took %.2f seconds to load the dictionary", elapsed.Seconds()))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		results := solver.GetValidAnagrams(q)
		w.Header().Set("Content-Type", "application/json")

		// Marshal the map to JSON
		jsonData, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})
	logger.Info(fmt.Sprintf("Running the server at %d port", PORT))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		logger.Error(fmt.Sprintf("Failed to start the server at %d port", PORT))
	}
}
