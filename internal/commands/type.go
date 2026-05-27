package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleType(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)

	if !exists {
		return protocol.SimpleString("none")
	}

	return protocol.SimpleString(string(value.Type))
}
