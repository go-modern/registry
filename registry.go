package registry

import (
	"fmt"
	"sync"
)

// ErrKeyExists is returned on TryPut when a key already exists in the registry.
var ErrKeyExists = fmt.Errorf("key already exists")

// Registry is a thread-safe map that allows only one value per key.
type Registry[K any, V any] struct {
	name string
	data sync.Map
}

// New creates a new registry with the given name.
func New[K any, V any](name string) *Registry[K, V] {
	return &Registry[K, V]{
		name: name,
	}
}

// Get returns the value for the given key, or false if the key does not exist.
func (r *Registry[K, V]) Get(key K) (*V, bool) {
	value, ok := r.data.Load(key)
	if !ok {
		return nil, false
	}
	return value.(*V), true
}

// Put sets the value for the given key. If the key already exists, it panics.
func (r *Registry[K, V]) Put(key K, value V) {
	if err := r.TryPut(key, value); err != nil {
		panic(err)
	}
}

// TryPut sets the value for the given key. If the key already exists, it returns an error.
func (r *Registry[K, V]) TryPut(key K, value V) error {
	if _, exists := r.data.LoadOrStore(key, &value); exists {
		return fmt.Errorf("registry(%v): %w: %v", r.name, ErrKeyExists, key)
	}
	return nil
}
