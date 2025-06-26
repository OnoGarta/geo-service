package auth

import "sync"

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]*User
}

func NewStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]*User)}
}

func (s *MemoryStore) Create(u *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[u.Username]; ok {
		return ErrExists
	}
	s.data[u.Username] = u
	return nil
}

func (s *MemoryStore) Get(username string) (*User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.data[username]
	return u, ok
}
