package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleSAdd(cl *client.Client, cmd []string) string {
	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if exists && value.Type != storage.SetType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	members, ok := db.Store.GetSet(key)
	if !ok {
		members = make(map[string]struct{})
	}

	added := 0
	for i := 2; i < len(cmd); i++ {
		member := cmd[i]
		if _, present := members[member]; !present {
			members[member] = struct{}{}
			added++
		}
	}

	if added > 0 || !exists {
		db.Store.SetSet(key, members, value.ExpiresAt)
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(added)
}
