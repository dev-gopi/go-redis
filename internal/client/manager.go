package client

import "sync"

type ClientManager struct {
	mu      sync.RWMutex
	clients map[string]*Client
}

var Manager = NewManager()

func NewManager() *ClientManager {

	return &ClientManager{
		clients: make(map[string]*Client),
	}
}

func (m *ClientManager) Add(cl *Client) {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[cl.ID] = cl
}

func (m *ClientManager) Remove(id string) {

	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.clients, id)
}

func (m *ClientManager) Count() int {

	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.clients)
}
