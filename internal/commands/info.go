package commands

import (
	"fmt"

	"github.com/dev-gopi/go-redis/internal/metrics"
	"github.com/dev-gopi/go-redis/internal/protocol"
)

func HandleInfo() string {

	info := fmt.Sprintf(`# Server
redis_version:7.0.0
redis_mode:standalone

# Clients
connected_clients:%d

# Memory
used_memory:%d

# Keyspace
db0:keys=%d
`,
		metrics.ConnectedClients(),
		metrics.UsedMemory(),
		metrics.TotalKeys(),
	)

	return protocol.BulkString(info)
}
