package commands

import (
	"strconv"
	"strings"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleSScan(cl *client.Client, cmd []string) string {
	// SSCAN key cursor [MATCH pattern] [COUNT count]
	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	cursor := 0
	count := 10
	pattern := "*"

	if c, err := strconv.Atoi(cmd[2]); err == nil {
		cursor = c
	}

	for i := 3; i < len(cmd); i++ {
		switch strings.ToUpper(cmd[i]) {
		case "MATCH":
			if i+1 < len(cmd) {
				pattern = cmd[i+1]
			}
		case "COUNT":
			if i+1 < len(cmd) {
				if v, err := strconv.Atoi(cmd[i+1]); err == nil {
					count = v
				}
			}
		}
	}

	db := storage.GetClientDB(cl)

	value, exists := db.Store.GetValue(key)
	if !exists || value.Type != storage.SetType {
		return "*2\r\n$1\r\n0\r\n*0\r\n"
	}

	members, ok := db.Store.GetSet(key)
	if !ok {
		return "*2\r\n$1\r\n0\r\n*0\r\n"
	}

	filtered := make([]string, 0, len(members))
	needle := strings.Trim(pattern, "*")
	for member := range members {
		if pattern == "*" || strings.Contains(member, needle) {
			filtered = append(filtered, member)
		}
	}

	total := len(filtered)
	if cursor >= total {
		return "*2\r\n$1\r\n0\r\n*0\r\n"
	}

	end := cursor + count
	if end > total {
		end = total
	}

	nextCursor := 0
	if end < total {
		nextCursor = end
	}

	resp := "*2\r\n"
	resp += "$" + strconv.Itoa(len(strconv.Itoa(nextCursor))) + "\r\n" + strconv.Itoa(nextCursor) + "\r\n"
	resp += "*" + strconv.Itoa(end-cursor) + "\r\n"

	for i := cursor; i < end; i++ {
		resp += protocol.BulkString(filtered[i])
	}

	return resp
}
