package commands

import (
	"strconv"
	"strings"

	"github.com/dev-gopi/go-redis/internal/client"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func HandleScan(
	cl *client.Client,
	cmd []string,
) string {

	/*
		SCAN cursor MATCH pattern COUNT count
	*/

	cursor := 0
	count := 10
	pattern := "*"
	typeFilter := ""

	if len(cmd) >= 2 {
		c, err := strconv.Atoi(cmd[1])
		if err == nil {
			cursor = c
		}
	}

	for i := 2; i < len(cmd); i++ {

		switch strings.ToUpper(cmd[i]) {

		case "MATCH":

			if i+1 < len(cmd) {
				pattern = cmd[i+1]
			}

		case "COUNT":

			if i+1 < len(cmd) {

				c, err := strconv.Atoi(cmd[i+1])
				if err == nil {
					count = c
				}
			}
		case "TYPE":

			if i+1 < len(cmd) {
				typeFilter = cmd[i+1]
			}
		}
	}

	db := storage.GetClientDB(cl)
	keys := db.Store.Keys()

	filtered := make([]string, 0)

	for _, key := range keys {

		// MATCH filtering
		if pattern != "*" {
			if !strings.Contains(key, strings.Trim(pattern, "*")) {
				continue
			}
		}

		// TYPE filtering (case-insensitive, substring-aware)
		if typeFilter != "" {
			val, ok := db.Store.GetValue(key)
			if !ok {
				continue
			}

			vtype := strings.ToLower(string(val.Type))
			tf := strings.ToLower(typeFilter)

			if !(strings.EqualFold(string(val.Type), typeFilter) || strings.Contains(vtype, tf) || strings.Contains(tf, vtype)) {
				continue
			}
		}

		filtered = append(filtered, key)
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

	resultKeys := filtered[cursor:end]

	resp := "*2\r\n"

	resp += "$" + strconv.Itoa(len(strconv.Itoa(nextCursor))) +
		"\r\n" +
		strconv.Itoa(nextCursor) +
		"\r\n"

	resp += "*" + strconv.Itoa(len(resultKeys)) + "\r\n"

	for _, key := range resultKeys {

		resp += "$" + strconv.Itoa(len(key)) +
			"\r\n" +
			key +
			"\r\n"
	}

	return resp
}
