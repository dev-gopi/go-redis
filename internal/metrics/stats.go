package metrics

import (
	"runtime"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func ConnectedClients() int {
	return client.Manager.Count()
}

func TotalKeys() int {
	db, _ := storage.Manager.GetDB(0)

	return db.Store.Size()
}

func UsedMemory() uint64 {

	var m runtime.MemStats

	runtime.ReadMemStats(&m)

	return m.Alloc
}
