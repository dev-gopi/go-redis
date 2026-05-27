package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleJSONDel(cl *client.Client, cmd []string) string {
	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)
	deleted := 0
	if db.Store.Del(key) {
		deleted = 1
	}

	if deleted > 0 {
		if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
			return protocol.Error("AOF write failed")
		}
		if wal.Manager != nil {
			_ = wal.Manager.Write(cl.SelectedDB, cmd)
		}
	}

	return protocol.Integer(deleted)
}
