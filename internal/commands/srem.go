package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleSRem(cl *client.Client, cmd []string) string {
	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.Integer(0)
	}

	if value.Type != storage.SetType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	members, ok := db.Store.GetSet(key)
	if !ok {
		return protocol.Integer(0)
	}

	removed := 0
	for i := 2; i < len(cmd); i++ {
		if _, present := members[cmd[i]]; present {
			delete(members, cmd[i])
			removed++
		}
	}

	if len(members) == 0 {
		db.Store.Del(key)
	} else {
		db.Store.SetSet(key, members, value.ExpiresAt)
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(removed)
}
