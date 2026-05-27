package commands

import (
	"strconv"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleGetRange(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) != 4 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]

	start, err := strconv.Atoi(cmd[2])
	if err != nil {
		return protocol.Error("invalid start")
	}

	end, err := strconv.Atoi(cmd[3])
	if err != nil {
		return protocol.Error("invalid end")
	}

	db := storage.GetClientDB(cl)
	value, exists := db.Store.Get(key)

	if !exists {
		return protocol.BulkString("")
	}

	length := len(value)

	// Handle negative indexes
	if start < 0 {
		start = length + start
	}

	if end < 0 {
		end = length + end
	}

	// Bounds safety
	if start < 0 {
		start = 0
	}

	if end >= length {
		end = length - 1
	}

	if start > end || start >= length {
		return protocol.BulkString("")
	}

	result := value[start : end+1]

	return protocol.BulkString(result)
}
