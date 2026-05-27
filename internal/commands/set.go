package commands

import (
	"strconv"
	"strings"
	"time"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleSet(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	value := cmd[2]

	nx := false
	xx := false

	var expiresAt time.Time

	for i := 3; i < len(cmd); i++ {

		switch strings.ToUpper(cmd[i]) {

		case "NX":
			nx = true

		case "XX":
			xx = true

		case "EX":

			if i+1 >= len(cmd) {
				return protocol.Error("EX requires seconds")
			}

			seconds, err := strconv.Atoi(cmd[i+1])
			if err != nil {
				return protocol.Error("invalid EX value")
			}

			expiresAt = time.Now().
				Add(time.Duration(seconds) * time.Second)

			i++

		case "PX":

			if i+1 >= len(cmd) {
				return protocol.Error("PX requires milliseconds")
			}

			ms, err := strconv.Atoi(cmd[i+1])
			if err != nil {
				return protocol.Error("invalid PX value")
			}

			expiresAt = time.Now().
				Add(time.Duration(ms) * time.Millisecond)

			i++
		}
	}

	db := storage.GetClientDB(cl)

	exists := db.Store.Exists(key)

	if nx && exists {
		return protocol.NullBulkString()
	}

	if xx && !exists {
		return protocol.NullBulkString()
	}

	db.Store.Set(
		key,
		value,
		expiresAt,
	)

	err := aof.Manager.Write(cl.SelectedDB, cmd)

	if err != nil {
		return protocol.Error("AOF write failed")
	}

	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.SimpleString("OK")
}
