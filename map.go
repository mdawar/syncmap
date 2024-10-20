// Package syncmap provides a simple and generic map that is safe for concurrent use.
package syncmap

import "sync"

// Map is a map that is safe for concurrent use.
type Map[K comparable, V any] struct {
	mu sync.RWMutex

	// kv is the underlying map.
	kv map[K]V
}

// New creates a new [Map] that is safe for concurrent use by multiple goroutines.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		kv: make(map[K]V),
	}
}

// Len returns the length of the [Map].
func (m *Map[K, V]) Len() int {
	return len(m.kv)
}

// Get returns the value stored in the [Map] for a key, or the zero value if no
// value is present.
//
// The ok result indicates whether the value was found in the map.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	value, ok = m.kv[key]
	return value, ok
}

// Set sets the value for a key in the [Map].
func (m *Map[K, V]) Set(key K, value V) {
	m.kv[key] = value
}

// Delete deletes the value for a key in the [Map].
func (m *Map[K, V]) Delete(key K) {
	delete(m.kv, key)
}

// Clear deletes all the entries in the [Map].
func (m *Map[K, V]) Clear() {
	clear(m.kv)
}
