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

// Clear deletes all the entries in the [Map].
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	clear(m.kv)
	m.mu.Unlock()
}
