package commands

import (
	"strconv"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleZAdd(cl *client.Client, cmd []string) string {
	if len(cmd) < 4 || len(cmd)%2 != 0 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if exists && value.Type != storage.ZSetType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	members, ok := db.Store.GetZSet(key)
	if !ok {
		members = make(map[string]float64)
	}

	added := 0
	for i := 2; i < len(cmd); i += 2 {
		score, err := strconv.ParseFloat(cmd[i], 64)
		if err != nil {
			return protocol.Error("value is not a valid float")
		}

		member := cmd[i+1]
		if _, present := members[member]; !present {
			added++
		}
		members[member] = score
	}

	if added > 0 || !exists {
		db.Store.SetZSet(key, members, value.ExpiresAt)
	} else {
		// ensure updated scores are persisted even when no new members were added
		db.Store.SetZSet(key, members, value.ExpiresAt)
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(added)
}
