package registry

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	// ErrKeyExists is returned when a key already exists in the registry.
	ErrKeyExists = fmt.Errorf("key already exists")
	// ErrKeyNoExist is returned when a key does not exist in the registry.
	ErrKeyNoExist = fmt.Errorf("key does not exist")
)

// Registry is a thread-safe map that allows only one value per key.
type Registry[K any, V any] struct {
	name string
	data sync.Map

	defaultKey atomic.Pointer[K]
}

// New creates a new registry with the given name.
func New[K any, V any](name string) *Registry[K, V] {
	return &Registry[K, V]{
		name: name,
	}
}

// Init sets the default key for the registry. It returns the default vallue or an error if the key
// does not exist or is already set.
func (r *Registry[K, V]) Init(key K) (*V, error) {
	value := r.Load(key)
	if value == nil {
		return nil, fmt.Errorf("registry(%s): %w: %v", r.name, ErrKeyNoExist, key)
	}
	if !r.defaultKey.CompareAndSwap(nil, &key) {
		return nil, fmt.Errorf("registry(%s): %w: %v", r.name, ErrKeyExists, key)
	}
	return value, nil
}

// Default returns a pointer to the default key and value, or nil if the default is not set.
func (r *Registry[K, V]) Default() (*K, *V) {
	key := r.defaultKey.Load()
	if key == nil {
		return nil, nil
	}
	return key, r.Load(*key)
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
