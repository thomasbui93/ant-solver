package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"go.uber.org/zap"
)

const PORT = 3000

type HTTPServer struct {
	solver *Solver
	logger *zap.Logger
	cache  *expirable.LRU[string, map[int][]string]
}

func NewHTTPServer(logger *zap.Logger, solver *Solver) *HTTPServer {
	cache := expirable.NewLRU[string, map[int][]string](1000, nil, time.Minute)
	return &HTTPServer{
		solver: solver,
		logger: logger,
		cache:  cache,
	}
}

func (h *HTTPServer) Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		cacheKey, err := calHash(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results, ok := h.cache.Get(cacheKey)
		if !ok {
			results = h.solver.GetValidAnagrams(q)
			h.cache.Add(cacheKey, results)
		}

		w.Header().Set("Content-Type", "application/json")

		// Marshal the map to JSON
		jsonData, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	h.logger.Info(fmt.Sprintf("Running the server at %d port", PORT))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to start the server at %d port", PORT))
	}
}
