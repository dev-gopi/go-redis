package snapshot

import (
	"encoding/json"
	"os"

	"github.com/dev-gopi/go-redis/internal/logger"

	"github.com/dev-gopi/go-redis/internal/storage"
)

func Save(path string) error {
	data := snapshotData{
		Databases: make(map[int]map[string]storage.Value),
	}

	for _, db := range storage.AllDBs() {
		data.Databases[db.ID] = db.Store.Export()
	}

	bytes, err := json.MarshalIndent(
		data,
		"",
		"  ",
	)

	if err != nil {
		return err
	}

	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return err
	}

	logger.InfoLogger.Printf("Snapshot saved: %s (bytes=%d)", path, len(bytes))

	return nil
}
