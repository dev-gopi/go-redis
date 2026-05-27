package commands

import (
	"fmt"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleHSet(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) < 4 || len(cmd)%2 != 0 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	var hash map[string]string

	value, exists := db.Store.GetValue(key)

	if exists {

		if value.Type != storage.HashType {
			return protocol.Error(
				"WRONGTYPE Operation against a key holding the wrong kind of value",
			)
		}

		// value.Data may be stored as map[string]string or map[string]interface{}
		switch v := value.Data.(type) {
		case map[string]string:
			hash = v
		case map[string]any:
			hash = make(map[string]string, len(v))
			for kk, vv := range v {
				if s, ok := vv.(string); ok {
					hash[kk] = s
				} else {
					// fallback to string conversion for non-string values
					hash[kk] = fmt.Sprint(vv)
				}
			}
		default:
			// unexpected underlying type
			return protocol.Error(
				"WRONGTYPE Operation against a key holding the wrong kind of value",
			)
		}

	} else {

		hash = make(map[string]string)
	}

	added := 0

	for i := 2; i < len(cmd); i += 2 {

		field := cmd[i]
		val := cmd[i+1]

		_, exists := hash[field]

		if !exists {
			added++
		}

		hash[field] = val
	}

	db.Store.SetValue(
		key,
		storage.Value{
			Type: storage.HashType,
			Data: hash,
		},
	)

	return protocol.Integer(added)
}
