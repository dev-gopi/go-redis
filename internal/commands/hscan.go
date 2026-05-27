package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/protocol"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleHScan(
	cl *client.Client,
	cmd []string,
) string {

	// HSCAN key cursor [MATCH pattern] [COUNT count]

	if len(cmd) < 3 {
		return protocol.Error("wrong number of arguments")
	}

	key := cmd[1]
	cursor := 0
	count := 10
	pattern := "*"

	c, err := strconv.Atoi(cmd[2])
	if err == nil {
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
	if !exists || value.Type != storage.HashType {
		// Return cursor 0 and empty array
		return "*2\r\n$1\r\n0\r\n*0\r\n"
	}

	// Extract fields
	var fields []string
	var values []string

	switch h := value.Data.(type) {
	case map[string]string:
		for fk, fv := range h {
			fields = append(fields, fk)
			values = append(values, fv)
		}
	case map[string]any:
		for fk, fv := range h {
			fields = append(fields, fk)
			values = append(values, fmt.Sprint(fv))
		}
	default:
		// unknown underlying type
		return protocol.Error("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	// filter by pattern
	filteredFields := make([]string, 0)
	filteredValues := make([]string, 0)

	for i, f := range fields {
		if pattern == "*" || strings.Contains(f, strings.Trim(pattern, "*")) {
			filteredFields = append(filteredFields, f)
			filteredValues = append(filteredValues, values[i])
		}
	}

	total := len(filteredFields)
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

	// cursor as bulk string
	resp += "$" + strconv.Itoa(len(strconv.Itoa(nextCursor))) + "\r\n" + strconv.Itoa(nextCursor) + "\r\n"

	// inner array of field,value pairs
	pairCount := (end - cursor) * 2
	resp += "*" + strconv.Itoa(pairCount) + "\r\n"

	for i := cursor; i < end; i++ {
		f := filteredFields[i]
		v := filteredValues[i]
		resp += protocol.BulkString(f)
		resp += protocol.BulkString(v)
	}

	return resp
}
