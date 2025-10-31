package internal

import (
	"fmt"
	"os"
	"path/filepath"
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
	absPath, err := filepath.Abs(s.folderPath)
	files, err := os.ReadDir(absPath)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	localMaps := make([]map[string][]string, len(files))
	for i, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			filePath := filepath.Join(s.folderPath, file.Name())
			wi := NewWordsImporter(s.logger, filePath)
			localMap, err := wi.ReadWords()
			if err != nil {
				s.logger.Error("Failed to read words from file", zap.String("file", file.Name()))
			}
			localMaps[i] = localMap
		}(file)
	}

	// wait for files processed
	wg.Wait()
	for _, lm := range localMaps {
		for hash, words := range lm {
			s.hd[hash] = append(s.hd[hash], words...)
		}
	}
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
