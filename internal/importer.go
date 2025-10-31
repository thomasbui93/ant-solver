package internal

import (
	"bufio"
	"io"
	"os"

	"go.uber.org/zap"
)

type WordsImporter struct {
	filePath string
	logger   *zap.Logger
}

func NewWordsImporter(logger *zap.Logger, filePath string) *WordsImporter {
	return &WordsImporter{
		filePath: filePath,
		logger:   logger,
	}
}

func (w *WordsImporter) ReadWords() (map[string][]string, error) {
	fs, err := os.OpenFile(w.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer fs.Close()
	dict := make(map[string][]string)
	buff := bufio.NewReader(fs)
	for {
		b, _, err := buff.ReadLine()
		if err != nil {
			if io.EOF == err {
				return dict, nil
			}
			return nil, err
		}
		word := string(b)
		if len(word) >= 3 {
			hash, err := calHash(word)
			if err != nil {
				continue
			}
			val, exist := dict[hash]
			if exist {
				dict[hash] = append(val, word)
			} else {
				dict[hash] = []string{word}
			}
		}
	}
}
