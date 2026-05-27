package wal

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dev-gopi/go-redis/internal/logger"
)

type Entry struct {
	Command []string
}

type Writer struct {
	mu   sync.Mutex
	file *os.File
}

var Manager *Writer

func Init(path string) error {

	file, err := os.OpenFile(
		path,
		os.O_CREATE|
			os.O_APPEND|
			os.O_RDWR,
		0644,
	)

	if err != nil {
		return err
	}

	Manager = &Writer{file: file}

	logger.InfoLogger.Printf("WAL initialized: %s", path)

	return nil
}

func (w *Writer) Write(dbID int, command []string) error {

	w.mu.Lock()
	defer w.mu.Unlock()

	if dbID != 0 {
		_, err := w.file.WriteString(fmt.Sprintf("SELECT %d\n", dbID))
		if err != nil {
			return err
		}
	}

	logger.InfoLogger.Printf("WAL write db=%d cmd=%v", dbID, command)

	_, err := w.file.WriteString(strings.Join(command, " ") + "\n")
	if err != nil {
		return err
	}

	return w.file.Sync()
}
