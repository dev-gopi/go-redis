package commands

import (
	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleDBSize(
	cl *client.Client,
	cmd []string,
) string {

	db := storage.GetClientDB(cl)

	size := db.Store.Size()

	return protocol.Integer(size)
}
