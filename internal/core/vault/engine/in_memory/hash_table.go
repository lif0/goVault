package in_memory

import "sync"

type HashTable struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
		mu:   sync.RWMutex{},
	}
}

func (s *HashTable) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *HashTable) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, found := s.data[key]
	return value, found
}

func (s *HashTable) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, key)
}
