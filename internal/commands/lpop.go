package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleLPop(cl *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	items, exists := db.Store.GetList(key)
	if !exists {
		value, ok := db.Store.GetValue(key)
		if ok && value.Type != storage.ListType {
			return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return protocol.NullBulkString()
	}

	if len(items) == 0 {
		db.Store.Del(key)
		return protocol.NullBulkString()
	}

	removed := items[0]
	remaining := items[1:]

	if len(remaining) == 0 {
		db.Store.Del(key)
	} else if value, ok := db.Store.GetValue(key); ok {
		db.Store.SetValue(key, storage.Value{Type: storage.ListType, Data: storage.ListValue(remaining), ExpiresAt: value.ExpiresAt})
	} else {
		db.Store.SetList(key, remaining, storage.Value{}.ExpiresAt)
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.BulkString(removed)
}
