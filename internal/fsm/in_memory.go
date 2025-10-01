package fsm

import (
	"errors"
	"sync"
)

type InMemoryStateStore struct {
	mu    sync.RWMutex
	store map[int64]string
}

func NewInMemoryStateStore() *InMemoryStateStore {
	return &InMemoryStateStore{store: make(map[int64]string)}
}

func (s *InMemoryStateStore) Set(userID int64, state string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[userID] = state
	return nil
}

func (s *InMemoryStateStore) Get(userID int64) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.store[userID]
	if !ok {
		return "", errors.New("cannot get state")
	}
	return st, nil
}

func (s *InMemoryStateStore) Delete(userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, userID)
	return nil
}
