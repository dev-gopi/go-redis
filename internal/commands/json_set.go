package commands

import (
	"encoding/json"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/persistence/aof"
	"github.com/dev-gopi/go-redis/internal/persistence/wal"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleJSONSet(cl *client.Client, cmd []string) string {
	if len(cmd) < 4 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	path := cmd[2]
	raw := cmd[3]

	if path != "." {
		return protocol.Error("only root path is supported")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		return protocol.Error("invalid JSON")
	}

	db := storage.GetClientDB(cl)
	db.Store.SetValue(key, storage.Value{Type: storage.JsonType, Data: storage.JSONValue{Raw: []byte(raw), Parsed: parsed}})

	if err := aof.Manager.Write(cl.SelectedDB, cmd); err != nil {
		return protocol.Error("AOF write failed")
	}
	if wal.Manager != nil {
		_ = wal.Manager.Write(cl.SelectedDB, cmd)
	}

	return protocol.SimpleString("OK")
}
