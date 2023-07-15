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

// Load returns a pointer to the value for the given key, or nil if the key does not exist.
func (r *Registry[K, V]) Load(key K) *V {
	value, ok := r.data.Load(key)
	if !ok {
		return nil
	}
	return value.(*V)
}

// MustStore sets the value for the given key. If the key already exists, it panics.
func (r *Registry[K, V]) MustStore(key K, value V) {
	if err := r.Store(key, value); err != nil {
		panic(err)
	}
}

// Store sets the value for the given key. If the key already exists, it returns an error.
func (r *Registry[K, V]) Store(key K, value V) error {
	if _, exists := r.data.LoadOrStore(key, &value); exists {
		return fmt.Errorf("registry(%s): %w: %v", r.name, ErrKeyExists, key)
	}
	return nil
}
