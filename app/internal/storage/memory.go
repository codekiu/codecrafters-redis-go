package storage

import (
	"sync"
	"time"
)

// MemoryStorage implements Storage using an in-memory map
type MemoryStorage struct {
	data  map[string]string
	mutex sync.RWMutex
}

// NewMemoryStorage creates a new memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

// Set stores a value for a key
func (s *MemoryStorage) Set(key string, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
}

// SetWithExpiry stores a value with an expiration time
func (s *MemoryStorage) SetWithExpiry(key string, value string, expiry time.Duration) {
	s.Set(key, value)

	go func() {
		<-time.After(expiry)
		s.Delete(key)
	}()
}

// Delete removes a key from storage
func (s *MemoryStorage) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
}

func (s *MemoryStorage) FlushAll() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for k := range s.data {
		delete(s.data, k)
	}
}
