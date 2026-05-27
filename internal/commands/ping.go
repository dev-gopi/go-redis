package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
)

func HandlePing(cl *client.Client, cmd []string) string {
	_ = cl

	if len(cmd) == 1 {
		return protocol.SimpleString("PONG")
	}

	if len(cmd) == 2 {
		return protocol.BulkString(cmd[1])
	}

	return protocol.Error("wrong number of arguments")
}
