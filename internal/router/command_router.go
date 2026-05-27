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
		return commands.HandlePing(cl, cmd)

	case "AUTH":
		return commands.HandleAuth(cl, cmd)

	case "INFO":
		return commands.HandleInfo()

	case "CLIENT":
		return commands.HandleClient(cl, cmd)

	case "SET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleSet(cl, cmd)

	case "SELECT":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleSelect(cl, cmd)

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

	case "EXISTS":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleExists(cl, cmd)

	case "DBSIZE":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleDBSize(cl, cmd)

	case "SCAN":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleScan(cl, cmd)

	case "TYPE":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleType(cl, cmd)

	case "TTL":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleTTL(cl, cmd)

	case "MEMORY":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleMemory(cl, cmd)

	case "STRLEN":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleStrlen(cl, cmd)

	case "GETRANGE":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleGetRange(cl, cmd)

	case "JSON.SET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleJSONSet(cl, cmd)

	case "JSON.GET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleJSONGet(cl, cmd)

	case "JSON.DEL":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleJSONDel(cl, cmd)

	case "HSET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHSet(cl, cmd)

	case "HGET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHGet(cl, cmd)

	case "HEXISTS":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHExists(cl, cmd)

	case "HDEL":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHDel(cl, cmd)

	case "HMGET":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHMGet(cl, cmd)

	case "HLEN":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHLen(cl, cmd)

	case "HSCAN":
		authErr := RequireAuth(cl)
		if authErr != "" {
			return authErr
		}
		return commands.HandleHScan(cl, cmd)

	default:
		return protocol.Error("unknown command")
	}
}
