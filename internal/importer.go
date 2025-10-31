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
	ch       chan string
}

func NewWordsImporter(logger *zap.Logger, filePath string, ch chan string) *WordsImporter {
	return &WordsImporter{
		filePath: filePath,
		logger:   logger,
		ch:       ch,
	}
}

func (w *WordsImporter) ReadWords() error {
	fs, err := os.OpenFile(w.filePath, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fs.Close()

	buff := bufio.NewReader(fs)
	for {
		b, _, err := buff.ReadLine()
		if err != nil {
			if io.EOF == err {
				return nil
			}
			return err
		}
		word := string(b)
		if len(word) >= 3 {
			w.ch <- word
		}
	}
}
