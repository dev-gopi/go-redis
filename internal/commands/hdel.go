package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleHDel(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.Integer(0)
	}

	if value.Type != storage.HashType {
		return protocol.Error(
			"WRONGTYPE Operation against a key holding the wrong kind of value",
		)
	}

	hash, ok := value.Data.(map[string]string)
	if !ok {
		return protocol.Integer(0)
	}

	deleted := 0
	for _, field := range cmd[2:] {
		if _, ok := hash[field]; ok {
			delete(hash, field)
			deleted++
		}
	}

	if deleted > 0 {
		if len(hash) == 0 {
			db.Store.Del(key)
		} else {
			db.Store.SetValue(
				key,
				storage.Value{Type: storage.HashType, Data: hash},
			)
		}

		if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
			return protocol.Error("AOF write failed")
		}

		if wal.Manager != nil {
			_ = wal.Manager.Write(cl.SelectedDB, cmd)
		}
	}

	return protocol.Integer(deleted)
}
