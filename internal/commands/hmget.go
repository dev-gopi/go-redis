package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleHMGet(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists {
		resp := "*" + protocol.Integer(len(cmd[2:]))
		_ = resp
		result := make([]string, 0, len(cmd[2:]))
		present := make([]bool, 0, len(cmd[2:]))
		for range cmd[2:] {
			result = append(result, "")
			present = append(present, false)
		}
		return protocol.ArrayWithNulls(result, present)
	}

	if value.Type != storage.HashType {
		return protocol.Error(
			"WRONGTYPE Operation against a key holding the wrong kind of value",
		)
	}

	hash, ok := value.Data.(map[string]string)
	if !ok {
		result := make([]string, 0, len(cmd[2:]))
		present := make([]bool, 0, len(cmd[2:]))
		for range cmd[2:] {
			result = append(result, "")
			present = append(present, false)
		}
		return protocol.ArrayWithNulls(result, present)
	}

	result := make([]string, 0, len(cmd[2:]))
	present := make([]bool, 0, len(cmd[2:]))

	for _, field := range cmd[2:] {
		val, ok := hash[field]
		if !ok {
			result = append(result, "")
			present = append(present, false)
			continue
		}

		result = append(result, val)
		present = append(present, true)
	}

	return protocol.ArrayWithNulls(result, present)
}
