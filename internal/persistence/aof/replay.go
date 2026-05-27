package aof

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
	"time"

	"strconv"

	"github.com/dev-gopi/go-redis/internal/logger"
	"github.com/dev-gopi/go-redis/internal/persistence/snapshot"
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

	// if a snapshot was loaded and the aof file is older-or-equal to the snapshot,
	// skip replay to avoid duplicating state already present in the snapshot.
	if !snapshot.LoadedAt.IsZero() {
		if fi, err := os.Stat(path); err == nil {
			if !fi.ModTime().After(snapshot.LoadedAt) {
				logger.InfoLogger.Printf("AOF is older-or-equal to snapshot (aof=%v snapshot=%v), skipping replay", fi.ModTime(), snapshot.LoadedAt)
				return nil
			}
		}
	}

	scanner := bufio.NewScanner(file)
	currentDB := 0

	for scanner.Scan() {

		line := scanner.Text()

		// Try JSON first (new format)
		var entry struct {
			DB  int      `json:"db"`
			Cmd []string `json:"cmd"`
		}

		parts := []string{}
		command := ""

		if err := json.Unmarshal([]byte(line), &entry); err == nil && len(entry.Cmd) > 0 {
			// JSON format
			currentDB = entry.DB
			parts = entry.Cmd
			command = strings.ToUpper(parts[0])
		} else {
			// Fallback to legacy space-separated format
			parts = strings.Split(line, " ")
			if len(parts) == 0 {
				continue
			}
			command = strings.ToUpper(parts[0])
		}

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

		case "LPUSH":

			if len(parts) < 3 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			items, exists := db.Store.GetList(parts[1])
			if !exists {
				value, ok := db.Store.GetValue(parts[1])
				if ok && value.Type != storage.ListType {
					continue
				}
			}

			for i := 2; i < len(parts); i++ {
				items = append([]string{parts[i]}, items...)
			}

			if value, ok := db.Store.GetValue(parts[1]); ok {
				db.Store.SetValue(parts[1], storage.Value{Type: storage.ListType, Data: storage.ListValue(items), ExpiresAt: value.ExpiresAt})
			} else {
				db.Store.SetList(parts[1], items, time.Time{})
			}

		case "LPUSHX":

			if len(parts) < 3 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(parts[1])
			if !exists || value.Type != storage.ListType {
				continue
			}

			items, ok := db.Store.GetList(parts[1])
			if !ok {
				continue
			}

			for i := 2; i < len(parts); i++ {
				items = append([]string{parts[i]}, items...)
			}

			db.Store.SetValue(parts[1], storage.Value{Type: storage.ListType, Data: storage.ListValue(items), ExpiresAt: value.ExpiresAt})

		case "RPUSH":

			if len(parts) < 3 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			items, exists := db.Store.GetList(parts[1])
			if !exists {
				value, ok := db.Store.GetValue(parts[1])
				if ok && value.Type != storage.ListType {
					continue
				}
			}

			items = append(items, parts[2:]...)

			if value, ok := db.Store.GetValue(parts[1]); ok {
				db.Store.SetValue(parts[1], storage.Value{Type: storage.ListType, Data: storage.ListValue(items), ExpiresAt: value.ExpiresAt})
			} else {
				db.Store.SetList(parts[1], items, time.Time{})
			}

		case "RPUSHX":

			if len(parts) < 3 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(parts[1])
			if !exists || value.Type != storage.ListType {
				continue
			}

			items, ok := db.Store.GetList(parts[1])
			if !ok {
				continue
			}

			items = append(items, parts[2:]...)
			db.Store.SetValue(parts[1], storage.Value{Type: storage.ListType, Data: storage.ListValue(items), ExpiresAt: value.ExpiresAt})

		case "LPOP":

			if len(parts) < 2 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			items, exists := db.Store.GetList(parts[1])
			if !exists || len(items) == 0 {
				db.Store.Del(parts[1])
				continue
			}

			remaining := items[1:]
			if len(remaining) == 0 {
				db.Store.Del(parts[1])
			} else if value, ok := db.Store.GetValue(parts[1]); ok {
				db.Store.SetValue(parts[1], storage.Value{Type: storage.ListType, Data: storage.ListValue(remaining), ExpiresAt: value.ExpiresAt})
			} else {
				db.Store.SetList(parts[1], remaining, time.Time{})
			}

		case "EXPIRE":

			if len(parts) != 3 {
				continue
			}

			seconds, err := strconv.Atoi(parts[2])
			if err != nil {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(parts[1])
			if !exists {
				continue
			}

			value.ExpiresAt = time.Now().Add(time.Duration(seconds) * time.Second)
			db.Store.SetValue(parts[1], value)

		case "INCR":

			if len(parts) != 2 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(parts[1])
			current := 0
			if exists {
				parsed, err := strconv.Atoi(value.Data.(string))
				if err != nil {
					continue
				}
				current = parsed
			}

			current++
			if exists {
				value.Data = strconv.Itoa(current)
				db.Store.SetValue(parts[1], value)
			} else {
				db.Store.Set(parts[1], strconv.Itoa(current), time.Time{})
			}

		case "SADD":

			if len(parts) < 3 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(parts[1])
			if exists && value.Type != storage.SetType {
				continue
			}

			members, ok := db.Store.GetSet(parts[1])
			if !ok {
				members = make(map[string]struct{})
			}

			for _, member := range parts[2:] {
				members[member] = struct{}{}
			}

			expiresAt := time.Time{}
			if exists {
				expiresAt = value.ExpiresAt
			}
			db.Store.SetSet(parts[1], members, expiresAt)

		case "SREM":

			if len(parts) < 3 {
				continue
			}

			db, err := storage.Manager.GetDB(currentDB)
			if err != nil {
				continue
			}

			value, exists := db.Store.GetValue(parts[1])
			if !exists || value.Type != storage.SetType {
				continue
			}

			members, ok := db.Store.GetSet(parts[1])
			if !ok {
				continue
			}

			for _, member := range parts[2:] {
				delete(members, member)
			}

			if len(members) == 0 {
				db.Store.Del(parts[1])
			} else {
				db.Store.SetSet(parts[1], members, value.ExpiresAt)
			}
		}
	}

	logger.InfoLogger.Println(
		"AOF replay completed",
	)

	return nil
}
