package commands

import "github.com/dev-gopi/go-redis/internal/protocol"

func HandleInfo() string {

	info := `
# Server
redis_version:7.0.0
redis_mode:standalone

# Clients
connected_clients:1

# Memory
used_memory:1024
`

	return protocol.BulkString(info)
}
