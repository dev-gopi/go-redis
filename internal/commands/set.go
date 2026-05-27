package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleSet(cl *client.Client, cmd []string) string {
	if len(cmd) != 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	value := cmd[2]

	storage.GlobalStore.Set(key, value)

	return protocol.SimpleString("OK")
}
