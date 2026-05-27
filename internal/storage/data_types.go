package storage

import (
	"sync"
	"time"
)

type ValueType string

const (
	StringType ValueType = "string"
	HashType   ValueType = "hash"
	ListType   ValueType = "list"
	SetType    ValueType = "set"
	ZSetType   ValueType = "zset"
	StreamType ValueType = "stream"
	JsonType   ValueType = "json"
)

type Value struct {
	Type      ValueType
	Data      any
	ExpiresAt time.Time
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Value
}
