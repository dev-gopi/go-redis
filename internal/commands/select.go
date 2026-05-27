package commands

import (
	"strconv"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleSelect(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	dbID, err := strconv.Atoi(cmd[1])
	if err != nil {
		return protocol.Error("invalid DB index")
	}

	_, err = storage.Manager.GetDB(dbID)
	if err != nil {
		return protocol.Error("DB index out of range")
	}

	cl.SelectedDB = dbID

	return protocol.SimpleString("OK")
}
