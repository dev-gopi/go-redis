package aof

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"

	"strconv"

	"github.com/dev-gopi/go-redis/internal/logger"
	"github.com/dev-gopi/go-redis/internal/storage"
)

func Replay(path string) error {

	file, err := os.Open(path)

	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentDB := 0

	for scanner.Scan() {

		line := scanner.Text()

		parts := strings.Split(line, " ")

		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])

		switch command {

		case "SELECT":

			if len(parts) < 2 {
				continue
			}

			dbID, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}

			if _, err := storage.Manager.GetDB(dbID); err != nil {
				continue
			}

			currentDB = dbID

		case "SET":

			if len(parts) < 3 {
				continue
			}

			key := parts[1]
			value := parts[2]

			var expiresAt time.Time

			// parse EX support
			for i := 3; i < len(parts); i++ {

				if strings.ToUpper(parts[i]) == "EX" {

					if i+1 < len(parts) {

						seconds := parts[i+1]

						ttl, err := time.ParseDuration(
							seconds + "s",
						)

						if err == nil {

							expiresAt = time.Now().
								Add(ttl)
						}
					}
				}
			}

			db, _ := storage.Manager.GetDB(currentDB)

			db.Store.SetValue(
				key,
				storage.Value{
					Type:      storage.StringType,
					Data:      value,
					ExpiresAt: expiresAt,
				},
			)

		case "HSET":

			if len(parts) < 4 || len(parts)%2 != 0 {
				continue
			}

			key := parts[1]
			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(key)
			var hash map[string]string

			if exists && value.Type == storage.HashType {
				hash, _ = value.Data.(map[string]string)
			}

			if hash == nil {
				hash = make(map[string]string)
			}

			for i := 2; i < len(parts); i += 2 {
				hash[parts[i]] = parts[i+1]
			}

			db.Store.SetValue(
				key,
				storage.Value{
					Type: storage.HashType,
					Data: hash,
				},
			)

		case "HDEL":

			if len(parts) < 3 {
				continue
			}

			key := parts[1]
			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(key)
			if !exists || value.Type != storage.HashType {
				continue
			}

			hash, ok := value.Data.(map[string]string)
			if !ok {
				continue
			}

			for _, field := range parts[2:] {
				delete(hash, field)
			}

			if len(hash) == 0 {
				db.Store.Del(key)
			} else {
				db.Store.SetValue(
					key,
					storage.Value{
						Type: storage.HashType,
						Data: hash,
					},
				)
			}

		case "JSON.SET":

			if len(parts) < 4 {
				continue
			}

			key := parts[1]
			path := parts[2]
			raw := parts[3]

			if path != "." {
				continue
			}

			var parsed map[string]any
			if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			db.Store.SetValue(
				key,
				storage.Value{
					Type: storage.JsonType,
					Data: storage.JSONValue{Raw: []byte(raw), Parsed: parsed},
				},
			)

		case "JSON.DEL":

			if len(parts) < 2 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			db.Store.Del(parts[1])

		case "DEL":

			if len(parts) < 2 {
				continue
			}

			db, _ := storage.Manager.GetDB(currentDB)
			db.Store.Del(parts[1])
		}
	}

	logger.InfoLogger.Println(
		"AOF replay completed",
	)

	return nil
}
