package router

import (
	"github.com/dev-gopi/go-redis/internal/auth"
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
)

func RequireAuth(cl *client.Client) string {

	if !auth.Manager.Enabled {
		return ""
	}

	if !cl.Authenticated {
		return protocol.Error("authentication required")
	}

	return ""
}
