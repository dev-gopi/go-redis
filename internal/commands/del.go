package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleDel(cl *client.Client, cmd []string) string {
	if len(cmd) != 2 {
		return protocol.Error("wrong number of arguments")
	}

	ok := storage.GlobalStore.Del(cmd[1])

	if ok {
		return protocol.Integer(1)
	}

	return protocol.Integer(0)
}
