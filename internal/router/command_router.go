package router

import (
	"strings"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/commands"
	"github.com/dev-gopi/go-redis/internal/protocol"
)

func Handle(cl *client.Client, cmd []string) string {

	if len(cmd) == 0 {
		return protocol.Error("empty command")
	}

	command := strings.ToUpper(cmd[0])

	switch command {

	case "PING":
		return protocol.SimpleString("PONG")

	case "SET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleSet(cl, cmd)

	case "GET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleGet(cl, cmd)

	case "DEL":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleDel(cl, cmd)

	case "AUTH":
		return commands.HandleAuth(cl, cmd)

	case "INFO":
		return commands.HandleInfo()

	default:
		return protocol.Error("unknown command")
	}
}
