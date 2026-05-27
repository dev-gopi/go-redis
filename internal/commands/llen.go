package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleLLen(cl *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	items, exists := db.Store.GetList(key)
	if !exists {
		value, ok := db.Store.GetValue(key)
		if ok && value.Type != storage.ListType {
			return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return protocol.Integer(0)
	}

	return protocol.Integer(len(items))
}
