package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleRPushX(cl *client.Client, cmd []string) string {
	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.Integer(0)
	}

	if value.Type != storage.ListType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	items, ok := db.Store.GetList(key)
	if !ok {
		return protocol.Integer(0)
	}

	items = append(items, cmd[2:]...)
	db.Store.SetValue(key, storage.Value{Type: storage.ListType, Data: storage.ListValue(items), ExpiresAt: value.ExpiresAt})

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(len(items))
}
