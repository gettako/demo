package services

import (
	"encoding/json"
	"errors"
)

// InMemoryStorage is a simple map-based implementation of contracts.KVStore.
// It is intended to be registered as a Lazy binding, loading data only when requested.
type InMemoryStorage struct {
	store map[string]string
}

// NewInMemoryStorage creates a new ephemeral KVStore.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		store: make(map[string]string),
	}
}

// Set serializes and stores a value under the given key.
func (s *InMemoryStorage) Set(key string, val any) error {
	b, err := json.Marshal(val)
	if err != nil {
		return err
	}
	s.store[key] = string(b)
	return nil
}

// Get retrieves a value by key and unmarshals it into the provided pointer.
func (s *InMemoryStorage) Get(key string, ptr any) error {
	val, ok := s.store[key]
	if !ok {
		return errors.New("key not found")
	}
	return json.Unmarshal([]byte(val), ptr)
}

// Has checks if a key exists in the store.
func (s *InMemoryStorage) Has(key string) bool {
	_, ok := s.store[key]
	return ok
}

// Delete removes a key from the store.
func (s *InMemoryStorage) Delete(key string) error {
	delete(s.store, key)
	return nil
}

// Close simulates flushing data and shutting down.
func (s *InMemoryStorage) Close() error {
	s.store = nil
	return nil
}
