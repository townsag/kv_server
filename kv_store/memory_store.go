package kv_store

import (
	"sync"
	"errors"
)

type MemoryStore struct {
	data map[string]interface{}
	mu sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]interface{}),
		// don't need to explicitly initialize the mutex because it is
		// automatically initialized to it's zero value. The zero value
		// of a mutex is an unlocked mutex
	}
}

func (ms *MemoryStore) Get(key string) (value interface{}, err error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	value, exists := ms.data[key]
	if !exists {	// use the happy-path-last coding pattern
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (ms *MemoryStore) Set(key string, value interface{}) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.data[key] = value
	return nil
}

func (ms *MemoryStore) Delete(key string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exists := ms.data[key]; !exists {
		return errors.New("key not found")
	}
	delete(ms.data, key)
	return nil
}