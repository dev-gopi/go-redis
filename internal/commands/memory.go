package commands

import (
	"unsafe"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleMemory(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	subcommand := cmd[1]

	if subcommand != "usage" &&
		subcommand != "USAGE" {
		return protocol.Error("unsupported MEMORY subcommand")
	}

	key := cmd[2]
	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)

	if !exists {
		return protocol.NullBulkString()
	}

	size := len(key)

	switch data := value.Data.(type) {

	case string:
		size += len(data)

	case []string:

		for _, item := range data {
			size += len(item)
		}

	case map[string]string:

		for k, v := range data {
			size += len(k)
			size += len(v)
		}

	case map[string]any:

		for k := range data {
			size += len(k)
		}

	case storage.JSONValue:

		size += len(data.Raw)

		for k := range data.Parsed {
			size += len(k)
		}
	}

	size += int(unsafe.Sizeof(value))

	return protocol.Integer(size)
}
