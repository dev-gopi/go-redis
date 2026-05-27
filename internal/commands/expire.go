package commands

import (
	"strconv"
	"time"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleExpire(cl *client.Client, cmd []string) string {
	if len(cmd) != 3 {
		return protocol.Error("wrong number of arguments")
	}

	seconds, err := strconv.Atoi(cmd[2])
	if err != nil {
		return protocol.Error("invalid expire time")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.Integer(0)
	}

	value.ExpiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
	db.Store.SetValue(key, value)

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.Integer(1)
}
