package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
		http.ServeFile(w, r, "./static/index.min.html")
	})

	http.HandleFunc("/api/unscramble", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		l := r.URL.Query().Get("length")
		s := r.URL.Query().Get("start")
		e := r.URL.Query().Get("end")

		if len(q) == 0 {
			http.Error(w, "empty query string is fobidden", http.StatusBadRequest)
		}
		q = strings.ToLower(q)

		cond := getCond(l, s, e)

		cacheKey, err := getCacheKey(q, cond)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results, ok := h.cache.Get(cacheKey)
		if !ok {
			results = h.solver.GetValidAnagramsAdvanced(q, cond)
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
