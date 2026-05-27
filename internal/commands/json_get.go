package commands

import (
	"encoding/json"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleJSONGet(cl *client.Client, cmd []string) string {
	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	path := cmd[2]
	if path != "." {
		return protocol.NullBulkString()
	}

	db := storage.GetClientDB(cl)
	value, exists := db.Store.GetValue(key)
	if !exists || value.Type != storage.JsonType {
		return protocol.NullBulkString()
	}

	jsonValue, ok := value.Data.(storage.JSONValue)
	if !ok {
		return protocol.NullBulkString()
	}

	if len(cmd) > 3 && cmd[3] == "NOESCAPE" {
		return protocol.BulkString(string(jsonValue.Raw))
	}

	bytes, err := json.Marshal(jsonValue.Parsed)
	if err != nil {
		return protocol.NullBulkString()
	}

	return protocol.BulkString(string(bytes))
}
