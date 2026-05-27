package snapshot

import (
	"time"

	"github.com/dev-gopi/go-redis/internal/storage"
)

type snapshotData struct {
	Databases map[int]map[string]storage.Value `json:"databases"`
}

// LoadedAt is the time the snapshot was loaded. If zero, no snapshot was loaded.
var LoadedAt = (time.Time{})
