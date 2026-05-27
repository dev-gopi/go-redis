package aof

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/dev-gopi/go-redis/internal/logger"
)

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

	Manager = &Writer{
		file: file,
	}

	return nil
}

type aofEntry struct {
	DB  int      `json:"db"`
	Cmd []string `json:"cmd"`
}

func (w *Writer) Write(
	dbID int,
	command []string,
) error {

	w.mu.Lock()
	defer w.mu.Unlock()

	entry := aofEntry{DB: dbID, Cmd: command}
	b, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	logger.InfoLogger.Printf("AOF write db=%d cmd=%v", dbID, command)

	if _, err := w.file.WriteString(string(b) + "\n"); err != nil {
		return err
	}

	return w.file.Sync()
}
