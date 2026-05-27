package storage

import (
	"time"

	"github.com/dev-gopi/go-redis/internal/logger"
)

type ListValue []string

type SetValue map[string]struct{}

type ZSetValue map[string]float64

func (s *Store) SetList(key string, items []string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copyItems := append([]string(nil), items...)
	s.data[key] = Value{
		Type:      ListType,
		Data:      ListValue(copyItems),
		ExpiresAt: expiresAt,
	}

	logger.DebugLogger.Printf("SETLIST key=%s len=%d expires=%v", key, len(copyItems), expiresAt)
}

func (s *Store) GetList(key string) ([]string, bool) {
	s.mu.RLock()
	val, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if !val.ExpiresAt.IsZero() && time.Now().After(val.ExpiresAt) {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, false
	}

	switch list := val.Data.(type) {
	case ListValue:
		return append([]string(nil), list...), true
	case []string:
		return append([]string(nil), list...), true
	default:
		return nil, false
	}
}

func (s *Store) SetSet(key string, members map[string]struct{}, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copyMembers := make(SetValue, len(members))
	for member := range members {
		copyMembers[member] = struct{}{}
	}

	s.data[key] = Value{
		Type:      SetType,
		Data:      copyMembers,
		ExpiresAt: expiresAt,
	}

	logger.DebugLogger.Printf("SETSET key=%s size=%d expires=%v", key, len(copyMembers), expiresAt)
}

func (s *Store) GetSet(key string) (map[string]struct{}, bool) {
	s.mu.RLock()
	val, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if !val.ExpiresAt.IsZero() && time.Now().After(val.ExpiresAt) {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, false
	}

	switch set := val.Data.(type) {
	case SetValue:
		copySet := make(map[string]struct{}, len(set))
		for member := range set {
			copySet[member] = struct{}{}
		}
		return copySet, true
	case map[string]struct{}:
		copySet := make(map[string]struct{}, len(set))
		for member := range set {
			copySet[member] = struct{}{}
		}
		return copySet, true
	default:
		return nil, false
	}
}

func (s *Store) SetZSet(key string, members map[string]float64, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	copyMembers := make(ZSetValue, len(members))
	for member, score := range members {
		copyMembers[member] = score
	}

	s.data[key] = Value{
		Type:      ZSetType,
		Data:      copyMembers,
		ExpiresAt: expiresAt,
	}

	logger.DebugLogger.Printf("SETZSET key=%s size=%d expires=%v", key, len(copyMembers), expiresAt)
}

func (s *Store) GetZSet(key string) (map[string]float64, bool) {
	s.mu.RLock()
	val, ok := s.data[key]
	s.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if !val.ExpiresAt.IsZero() && time.Now().After(val.ExpiresAt) {
		s.mu.Lock()
		delete(s.data, key)
		s.mu.Unlock()
		return nil, false
	}

	switch zset := val.Data.(type) {
	case ZSetValue:
		copyZSet := make(map[string]float64, len(zset))
		for member, score := range zset {
			copyZSet[member] = score
		}
		return copyZSet, true
	case map[string]float64:
		copyZSet := make(map[string]float64, len(zset))
		for member, score := range zset {
			copyZSet[member] = score
		}
		return copyZSet, true
	default:
		return nil, false
	}
}
