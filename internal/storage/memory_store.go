package storage

import (
	"time"

	"github.com/dev-gopi/go-redis/internal/logger"
)

func NewStore() *Store {
	return &Store{
		data: make(map[string]Value),
	}
}

func (s *Store) Set(
	key string,
	value string,
	expiresAt time.Time,
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = Value{
		Type:      StringType,
		Data:      value,
		ExpiresAt: expiresAt,
	}

	logger.DebugLogger.Printf("SET key=%s type=string expires=%v", key, expiresAt)
}

func (s *Store) Get(
	key string,
) (string, bool) {

	s.mu.RLock()

	val, ok := s.data[key]

	s.mu.RUnlock()

	if !ok {
		return "", false
	}

	if !val.ExpiresAt.IsZero() &&
		time.Now().After(val.ExpiresAt) {

		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()

		return "", false
	}

	str, ok := val.Data.(string)

	if !ok {
		return "", false
	}
	logger.DebugLogger.Printf("GET key=%s ok=%v", key, ok)
	return str, true
}

func (s *Store) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]

	if exists {
		delete(s.data, key)
	}
	logger.DebugLogger.Printf("DEL key=%s existed=%v", key, exists)
	return exists
}

func (s *Store) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

func (s *Store) Keys() []string {

	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}
	logger.DebugLogger.Printf("KEYS count=%d", len(keys))
	return keys
}

func (s *Store) Exists(key string) bool {

	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.data[key]

	if !exists {
		return false
	}

	if !val.ExpiresAt.IsZero() &&
		time.Now().After(val.ExpiresAt) {

		delete(s.data, key)

		return false
	}

	return true
}

func (s *Store) GetValue(
	key string,
) (Value, bool) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]

	logger.DebugLogger.Printf("GetValue key=%s ok=%v type=%v", key, ok, val.Type)

	return val, ok
}

func (s *Store) RemoveExpiredKeys() {

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	for key, value := range s.data {

		if !value.ExpiresAt.IsZero() &&
			now.After(value.ExpiresAt) {

			delete(s.data, key)
		}
	}
}

func (s *Store) SetValue(
	key string,
	value Value,
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	logger.DebugLogger.Printf("SetValue key=%s type=%v", key, value.Type)
}

func (s *Store) Export() map[string]Value {

	s.mu.Lock()
	defer s.mu.Unlock()

	copyData := make(map[string]Value)
	now := time.Now()

	for k, v := range s.data {
		if !v.ExpiresAt.IsZero() && now.After(v.ExpiresAt) {
			delete(s.data, k)
			continue
		}
		copyData[k] = v
	}
	logger.DebugLogger.Printf("Export count=%d", len(copyData))
	return copyData
}

func (s *Store) Import(
	data map[string]Value,
) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = data

	logger.InfoLogger.Printf("Imported store with %d keys", len(data))
}
