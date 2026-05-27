package snapshot

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dev-gopi/go-redis/internal/logger"

	"github.com/dev-gopi/go-redis/internal/storage"
)

func Load(path string) error {

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
		db.Store.Import(normalizeValues(legacy))
		return nil
	}

	for dbID, values := range data.Databases {
		db, err := storage.Manager.GetDB(dbID)
		if err != nil {
			continue
		}

		db.Store.Import(normalizeValues(values))
	}

	logger.InfoLogger.Printf("Snapshot loaded: %s", path)

	return nil
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
