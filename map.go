// Package syncmap provides a simple and generic map that is safe for concurrent use.
package syncmap

import "sync"

// Map is a map that is safe for concurrent use.
type Map[K comparable, V any] struct {
	mu sync.RWMutex
}

// New creates a new [Map] that is safe for concurrent use by multiple goroutines.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

// Len returns the length of the [Map].
func (m *Map[K, V]) Len() int {
	return 0
}

// Get returns the value stored in the [Map] for a key, or the zero value if no
// value is present.
//
// The ok result indicates whether the value was found in the map.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	var zero V
	return zero, false
}

// Set sets the value for a key in the [Map].
func (m *Map[K, V]) Set(key K, value V) {
}

// Delete deletes the value for a key in the [Map].
func (m *Map[K, V]) Delete(key K) {
}

// Clear deletes all the entries in the [Map].
func (m *Map[K, V]) Clear() {
}
