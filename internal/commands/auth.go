package commands

import (
	"github.com/dev-gopi/go-redis/internal/auth"
	network "github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
)

func HandleAuth(
	cl *network.Client,
	cmd []string,
) string {

	// AUTH password
	if len(cmd) == 2 {

		ok := auth.Manager.Authenticate(
			"default",
			cmd[1],
		)

		if ok {

			cl.Authenticated = true
			cl.Username = "default"

			return protocol.SimpleString("OK")
		}

		return protocol.Error("invalid password")
	}

	// AUTH username password
	if len(cmd) == 3 {

		ok := auth.Manager.Authenticate(
			cmd[1],
			cmd[2],
		)

		if ok {

			cl.Authenticated = true
			cl.Username = cmd[1]

			return protocol.SimpleString("OK")
		}

		return protocol.Error("invalid username-password")
	}

	return protocol.Error("wrong number of arguments")
}
