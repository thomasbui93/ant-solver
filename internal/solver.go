package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const ALPHABET = 26
const MAX_WORD_LEN = 7

type Solver struct {
	folderPath string
	logger     *zap.Logger
	hd         map[string][]string
	mu         sync.RWMutex
}

func NewSolver(logger *zap.Logger, folderPath string) *Solver {
	return &Solver{
		folderPath: folderPath,
		logger:     logger,
		hd:         make(map[string][]string),
	}
}

func (s *Solver) Load() error {
	files, err := os.ReadDir(s.folderPath)
	if err != nil {
		return err
	}

	ch := make(chan string)

	// number of workers to process words concurrently
	numWorkers := runtime.NumCPU() // or any number you like

	var consumers sync.WaitGroup
	for _ = range numWorkers {
		consumers.Go(func() {
			for word := range ch {
				s.processWord(word)
			}
		})
	}

	// --- PRODUCERS ---
	var producerWG sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() {
			producerWG.Add(1)
			go func(file os.DirEntry) {
				defer producerWG.Done()
				filePath := filepath.Join(s.folderPath, file.Name())
				wi := NewWordsImporter(s.logger, filePath, ch)
				if err := wi.ReadWords(); err != nil {
					s.logger.Error("Failed to read words from file", zap.String("file", file.Name()))
				}
			}(file)
		}
	}

	// close channel after all producers done
	go func() {
		producerWG.Wait()
		close(ch)
	}()

	// wait for consumers to finish
	consumers.Wait()

	return nil
}

func (s *Solver) GetValidAnagrams(word string) map[int][]string {
	n := len(word)
	results := make(map[int]map[string]bool) // intermediate map to avoid duplicates
	seen := make(map[string]bool)

	for length := 1; length <= MAX_WORD_LEN && length <= n; length++ {
		var combs [][]int
		combinations(n, length, 0, []int{}, &combs)

		for _, indices := range combs {
			var freq [ALPHABET]int
			for _, idx := range indices {
				c := word[idx]
				if c >= 'a' && c <= 'z' {
					freq[c-'a']++
				}
			}
			key := freqKey(freq)
			if seen[key] {
				continue
			}
			seen[key] = true

			if matches, ok := s.hd[key]; ok {
				if results[length] == nil {
					results[length] = make(map[string]bool)
				}
				for _, w := range matches {
					results[length][w] = true
				}
			}
		}
	}

	// Convert to map[int][]string
	final := make(map[int][]string)
	for length, wordsMap := range results {
		for w := range wordsMap {
			final[length] = append(final[length], w)
		}
	}
	return final
}

func combinations(n, k int, start int, curr []int, all *[][]int) {
	if len(curr) == k {
		comb := make([]int, k)
		copy(comb, curr)
		*all = append(*all, comb)
		return
	}
	for i := start; i < n; i++ {
		curr = append(curr, i)
		combinations(n, k, i+1, curr, all)
		curr = curr[:len(curr)-1]
	}
}

func (s *Solver) processWord(word string) error {
	hash, err := calHash(word)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	val, exist := s.hd[hash]
	if exist {
		s.hd[hash] = append(val, word)
	} else {
		s.hd[hash] = []string{word}
	}

	return nil
}

func calHash(word string) (string, error) {
	var freq [ALPHABET]int
	for _, c := range word {
		if c >= 'a' && c <= 'z' {
			freq[c-'a']++
		} else {
			return "", fmt.Errorf("Invalid word :%s", word)
		}
	}
	return freqKey(freq), nil
}

func freqKey(freq [ALPHABET]int) string {
	parts := make([]string, ALPHABET)
	for i, count := range freq {
		parts[i] = fmt.Sprintf("%d", count)
	}
	return strings.Join(parts, "#")
}
