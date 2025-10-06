// Package syncmap provides a simple and generic map that is safe for concurrent use.
package syncmap

import (
	"iter"
	"sync"
)

// Map is a map that is safe for concurrent use.
type Map[K comparable, V any] struct {
	kv map[K]V // kv is the underlying map.
	mu sync.RWMutex
}

// New creates a new [Map] that is safe for concurrent use by multiple goroutines.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		kv: make(map[K]V),
	}
}

// NewWithCapacity creates a new [Map] with the specified capacity hint.
//
// The capacity hint does not bound the map size, it will create a map with
// an initial space to hold the specified number of elements.
func NewWithCapacity[K comparable, V any](capacity int) *Map[K, V] {
	return &Map[K, V]{
		kv: make(map[K]V, capacity),
	}
}

// Len returns the length of the [Map].
func (m *Map[K, V]) Len() int {
	m.mu.RLock()
	l := len(m.kv)
	m.mu.RUnlock()
	return l
}

// Get returns the value stored in the [Map] for a key, or the zero value if no
// value is present.
//
// The ok result indicates whether the value was found in the map.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	m.mu.RLock()
	value, ok = m.kv[key]
	m.mu.RUnlock()
	return value, ok
}

// Set sets the value for a key in the [Map].
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	m.kv[key] = value
	m.mu.Unlock()
}

// Delete deletes the value for a key in the [Map].
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.kv, key)
	m.mu.Unlock()
}

// Contains reports whether key is present in the [Map].
func (m *Map[K, V]) Contains(key K) bool {
	_, ok := m.Get(key)
	return ok
}

// Clear deletes all the entries in the [Map].
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	clear(m.kv)
	m.mu.Unlock()
}

// All returns an iterator over key-value pairs from the [Map].
//
// Similar to the map type, the iteration order is not guaranteed.
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		for key, value := range m.kv {
			if !yield(key, value) {
				return
			}
		}
	}
}
