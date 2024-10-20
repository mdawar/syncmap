package syncmap_test

import (
	"fmt"
	"maps"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/mdawar/syncmap"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestMapSetGet(t *testing.T) {
	t.Parallel()

	cases := []struct {
		key   string
		value int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}

	m := syncmap.New[string, int]()

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v-%v", tc.key, tc.value), func(t *testing.T) {
			m.Set(tc.key, tc.value)

			got, ok := m.Get(tc.key)
			require.True(t, ok, "key does not exist")
			assert.Equal(t, got, tc.value)
		})
	}
}

func TestMapLen(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, string]()

	m.Set("a", "a")
	require.Equal(t, 1, m.Len())

	m.Set("b", "b")
	require.Equal(t, 2, m.Len())

	// Key already exists.
	m.Set("a", "1")
	require.Equal(t, 2, m.Len())
}

func TestMapDelete(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, int]()

	key := "key"

	m.Set(key, 1)
	_, ok := m.Get(key)
	require.True(t, ok, "key does not exist")
	require.Equal(t, 1, m.Len())

	m.Set("keep", 2)
	require.Equal(t, 2, m.Len())

	m.Delete(key)
	_, ok = m.Get(key)
	require.False(t, ok, "key was not deleted")
	require.Equal(t, 1, m.Len())
}

func TestMapClear(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, int]()

	cases := []struct {
		key   string
		value int
	}{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
	}

	for _, tc := range cases {
		m.Set(tc.key, tc.value)
	}

	require.Equal(t, len(cases), m.Len())
	m.Clear()
	require.Equal(t, 0, m.Len())

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			_, ok := m.Get(tc.key)
			assert.False(t, ok, "key was not deleted on clear")
		})
	}
}

func TestMapAll(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, int]()

	want := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	for k, v := range want {
		m.Set(k, v)
	}

	got := make(map[string]int)
	for k, v := range m.All() {
		got[k] = v
	}

	if !maps.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestMapAllPartialIteration(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, int]()

	want := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	for k, v := range want {
		m.Set(k, v)
	}

	var key string
	var value int
	for k, v := range m.All() {
		key = k
		value = v
		break
	}

	got, ok := m.Get(key)
	require.True(t, ok, "key does not exist")
	assert.Equal(t, value, got)
}

func TestMapSetGetConcurrent(t *testing.T) {
	t.Parallel()

	m := syncmap.New[int, int]()

	var wg sync.WaitGroup

	// Write.
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Set(i, i)
		}()
	}

	// Read.
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Get(i)
		}()
	}

	wg.Wait()
}

func TestMapLenConcurrent(t *testing.T) {
	t.Parallel()

	m := syncmap.New[int, int]()

	var wg sync.WaitGroup

	// Write.
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Set(i, i)
		}()
	}

	// Read.
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Len()
		}()
	}

	wg.Wait()
}

func TestMapDeleteConcurrent(t *testing.T) {
	t.Parallel()

	m := syncmap.New[int, int]()

	var wg sync.WaitGroup

	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Delete(i)
		}()
	}

	wg.Wait()
}

func TestMapClearConcurrent(t *testing.T) {
	t.Parallel()

	m := syncmap.New[int, int]()

	var wg sync.WaitGroup

	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Clear()
		}()
	}

	wg.Wait()
}

func TestMapAllConcurrent(t *testing.T) {
	t.Parallel()

	m := syncmap.New[int, int]()

	var wg sync.WaitGroup

	// Write.
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Set(i, i)
		}()
	}

	// Read.
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range m.All() {
			}
		}()
	}

	wg.Wait()
}
