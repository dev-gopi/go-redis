package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleRPop(cl *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.NullBulkString()
	}

	if value.Type != storage.ListType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	items, ok := db.Store.GetList(key)
	if !ok || len(items) == 0 {
		db.Store.Del(key)
		return protocol.NullBulkString()
	}

	removed := items[len(items)-1]
	remaining := items[:len(items)-1]
	if len(remaining) == 0 {
		db.Store.Del(key)
	} else {
		db.Store.SetValue(key, storage.Value{Type: storage.ListType, Data: storage.ListValue(remaining), ExpiresAt: value.ExpiresAt})
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.BulkString(removed)
}
