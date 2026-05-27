package snapshot

import "github.com/dev-gopi/go-redis/internal/storage"

type snapshotData struct {
	Databases map[int]map[string]storage.Value `json:"databases"`
}
