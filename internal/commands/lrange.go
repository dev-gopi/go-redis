package commands

import (
	"strconv"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleLRange(cl *client.Client, cmd []string) string {
	if len(cmd) != 4 {
		return protocol.Error("wrong number of arguments")
	}

	start, err := strconv.Atoi(cmd[2])
	if err != nil {
		return protocol.Error("invalid start index")
	}

	stop, err := strconv.Atoi(cmd[3])
	if err != nil {
		return protocol.Error("invalid stop index")
	}

	key := cmd[1]
	db := storage.GetClientDB(cl)

	items, exists := db.Store.GetList(key)
	if !exists {
		value, ok := db.Store.GetValue(key)
		if ok && value.Type != storage.ListType {
			return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		return protocol.Array([]string{})
	}

	if len(items) == 0 {
		return protocol.Array([]string{})
	}

	length := len(items)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start >= length || start > stop {
		return protocol.Array([]string{})
	}

	return protocol.Array(items[start : stop+1])
}
