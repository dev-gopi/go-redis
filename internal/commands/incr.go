package commands

import (
	"strconv"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleIncr(cl *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	current := 0

	if exists {
		if value.Type != storage.StringType {
			return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
		}

		parsed, err := strconv.Atoi(value.Data.(string))
		if err != nil {
			return protocol.Error("value is not an integer or out of range")
		}
		current = parsed
	}

	current++

	if exists {
		value.Data = strconv.Itoa(current)
		db.Store.SetValue(key, value)
	} else {
		db.Store.Set(key, strconv.Itoa(current), value.ExpiresAt)
	}

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(current)
}
