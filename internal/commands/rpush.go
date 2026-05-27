package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleRPush(cl *client.Client, cmd []string) string {
	if len(cmd) < 3 {
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
	}

	for i := 2; i < len(cmd); i++ {
		items = append(items, cmd[i])
	}

	if value, ok := db.Store.GetValue(key); ok {
		db.Store.SetValue(key, storage.Value{Type: storage.ListType, Data: storage.ListValue(items), ExpiresAt: value.ExpiresAt})
	} else {
		db.Store.SetList(key, items, storage.Value{}.ExpiresAt)
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(len(items))
}
