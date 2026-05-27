package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleStrlen(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.Get(key)

	if !exists {
		return protocol.Integer(0)
	}

	return protocol.Integer(len(value))
}
