package commands

import (
	"strings"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
)

func HandleClient(
	cl *client.Client,
	cmd []string,
) string {

	if len(cmd) < 2 {
		return protocol.Error("wrong number of arguments")
	}

	subcommand := strings.ToUpper(cmd[1])

	switch subcommand {

	case "SETNAME":

		if len(cmd) != 3 {
			return protocol.Error("wrong number of arguments")
		}

		cl.Name = cmd[2]

		return protocol.SimpleString("OK")

	case "GETNAME":

		if cl.Name == "" {
			return protocol.NullBulkString()
		}

		return protocol.BulkString(cl.Name)

	default:
		return protocol.Error("unsupported CLIENT subcommand")
	}
}
