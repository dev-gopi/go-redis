package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleZCard(cl *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		return protocol.Integer(0)
	}

	if value.Type != storage.ZSetType {
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	members, ok := db.Store.GetZSet(key)
	if !ok {
		return protocol.Integer(0)
	}

	return protocol.Integer(len(members))
}
