package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleGet(client *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]

	value, ok := storage.GlobalStore.Get(key)
	if !ok {
		return protocol.NullBulkString()
	}

	return protocol.BulkString(value)
}
