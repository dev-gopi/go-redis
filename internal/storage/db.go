package storage

import (
	"fmt"
	"sync"

	"github.com/dev-gopi/go-redis/internal/client"
)

const DefaultDBCount = 16

type DB struct {
	ID    int
	Store *Store
}

type DBManager struct {
	mu  sync.RWMutex
	dbs map[int]*DB
}

var Manager = NewDBManager(DefaultDBCount)

func NewDBManager(count int) *DBManager {

	manager := &DBManager{
		dbs: make(map[int]*DB),
	}

	for i := 0; i < count; i++ {

		manager.dbs[i] = &DB{
			ID:    i,
			Store: NewStore(),
		}
	}

	return manager
}

func (m *DBManager) GetDB(id int) (*DB, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	db, ok := m.dbs[id]

	if !ok {
		return nil, fmt.Errorf("DB does not exist")
	}

	return db, nil
}

func (m *DBManager) ForEachDB(fn func(*DB)) {

	m.mu.RLock()
	dbs := make([]*DB, 0, len(m.dbs))
	for _, db := range m.dbs {
		dbs = append(dbs, db)
	}
	m.mu.RUnlock()

	for _, db := range dbs {
		fn(db)
	}
}

func GetClientDB(cl *client.Client) *DB {

	db, err := Manager.GetDB(cl.SelectedDB)
	if err != nil {
		db, _ = Manager.GetDB(0)
	}

	return db
}

func AllDBs() []*DB {

	Manager.mu.RLock()
	defer Manager.mu.RUnlock()

	dbs := make([]*DB, 0, len(Manager.dbs))
	for _, db := range Manager.dbs {
		dbs = append(dbs, db)
	}

	return dbs
}
