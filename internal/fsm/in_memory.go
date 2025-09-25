package fsm

import (
	"sync"
)

type InMemoryStateStore struct {
	mu    sync.RWMutex
	store map[int64]string
}

func NewInMemoryStateStore() *InMemoryStateStore {
	return &InMemoryStateStore{store: make(map[int64]string)}
}

func (s *InMemoryStateStore) Set(userID int64, state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[userID] = state
}

func (s *InMemoryStateStore) Get(userID int64) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.store[userID]
	return st, ok
}

func (s *InMemoryStateStore) Delete(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, userID)
}
