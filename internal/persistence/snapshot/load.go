package snapshot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dev-gopi/go-redis/internal/logger"

	"github.com/dev-gopi/go-redis/internal/storage"
)

func Load(path string) error {

	// capture file mod time
	var modTime time.Time
	if fi, err := os.Stat(path); err == nil {
		modTime = fi.ModTime()
	}

	bytes, err := os.ReadFile(path)

	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}

		if len(bytes) == 0 {
			return nil
		}

		return err
	}

	if len(bytes) == 0 {
		return nil
	}

	var data snapshotData

	err = json.Unmarshal(bytes, &data)

	if err != nil {
		var legacy map[string]storage.Value
		err = json.Unmarshal(bytes, &legacy)
		if err != nil {
			return err
		}

		db, _ := storage.Manager.GetDB(0)
		db.Store.Import(filterExpiredValues(normalizeValues(legacy)))
		return nil
	}

	for dbID, values := range data.Databases {
		db, err := storage.Manager.GetDB(dbID)
		if err != nil {
			continue
		}

		db.Store.Import(filterExpiredValues(normalizeValues(values)))
	}

	// set LoadedAt to snapshot file mod time
	if !modTime.IsZero() {
		LoadedAt = modTime
	} else {
		LoadedAt = time.Now()
	}

	logger.InfoLogger.Printf("Snapshot loaded: %s", path)

	return nil
}

func filterExpiredValues(values map[string]storage.Value) map[string]storage.Value {
	filtered := make(map[string]storage.Value, len(values))
	now := time.Now()

	for key, value := range values {
		if !value.ExpiresAt.IsZero() && now.After(value.ExpiresAt) {
			continue
		}
		filtered[key] = value
	}

	return filtered
}

func normalizeValues(values map[string]storage.Value) map[string]storage.Value {
	normalized := make(map[string]storage.Value, len(values))

	for key, value := range values {
		normalized[key] = normalizeValue(value)
	}

	return normalized
}

func normalizeValue(value storage.Value) storage.Value {
	switch value.Type {
	case storage.HashType:
		value.Data = normalizeHashData(value.Data)
	case storage.ListType:
		value.Data = normalizeListData(value.Data)
	case storage.SetType:
		value.Data = normalizeSetData(value.Data)
	case storage.ZSetType:
		value.Data = normalizeZSetData(value.Data)
	case storage.JsonType:
		value.Data = normalizeJSONData(value.Data)
	}

	return value
}

func normalizeHashData(data any) any {
	switch hash := data.(type) {
	case map[string]string:
		return hash
	case map[string]any:
		normalized := make(map[string]string, len(hash))
		for field, raw := range hash {
			normalized[field] = fmt.Sprint(raw)
		}
		return normalized
	default:
		return data
	}
}

func normalizeJSONData(data any) any {
	jsonValue, ok := data.(map[string]any)
	if !ok {
		return data
	}

	raw := []byte{}
	if rawText, ok := jsonValue["Raw"].(string); ok {
		if decoded, err := base64.StdEncoding.DecodeString(rawText); err == nil {
			raw = decoded
		} else {
			raw = []byte(rawText)
		}
	}

	parsed := make(map[string]any)
	if parsedValue, ok := jsonValue["Parsed"].(map[string]any); ok {
		parsed = parsedValue
	}

	return storage.JSONValue{Raw: raw, Parsed: parsed}
}

func normalizeListData(data any) any {
	switch list := data.(type) {
	case storage.ListValue:
		return list
	case []string:
		return storage.ListValue(list)
	case []any:
		normalized := make(storage.ListValue, 0, len(list))
		for _, item := range list {
			normalized = append(normalized, fmt.Sprint(item))
		}
		return normalized
	default:
		return data
	}
}

func normalizeSetData(data any) any {
	switch set := data.(type) {
	case storage.SetValue:
		return set
	case map[string]struct{}:
		copySet := make(storage.SetValue, len(set))
		for member := range set {
			copySet[member] = struct{}{}
		}
		return copySet
	case map[string]any:
		copySet := make(storage.SetValue, len(set))
		for member := range set {
			copySet[member] = struct{}{}
		}
		return copySet
	case []any:
		copySet := make(storage.SetValue, len(set))
		for _, member := range set {
			copySet[fmt.Sprint(member)] = struct{}{}
		}
		return copySet
	default:
		return data
	}
}

func normalizeZSetData(data any) any {
	switch zset := data.(type) {
	case storage.ZSetValue:
		return zset
	case map[string]float64:
		copyZSet := make(storage.ZSetValue, len(zset))
		for member, score := range zset {
			copyZSet[member] = score
		}
		return copyZSet
	case map[string]any:
		copyZSet := make(storage.ZSetValue, len(zset))
		for member, rawScore := range zset {
			switch score := rawScore.(type) {
			case float64:
				copyZSet[member] = score
			case int:
				copyZSet[member] = float64(score)
			case json.Number:
				if parsed, err := score.Float64(); err == nil {
					copyZSet[member] = parsed
				}
			default:
				if parsed, err := strconv.ParseFloat(fmt.Sprint(score), 64); err == nil {
					copyZSet[member] = parsed
				}
			}
		}
		return copyZSet
	default:
		return data
	}
}
