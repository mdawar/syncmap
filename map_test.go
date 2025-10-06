package syncmap_test

import (
	"fmt"
	"maps"
	"sync"
	"testing"

	"github.com/mdawar/syncmap"
)

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
			if !ok {
				t.Fatalf("key does not exist: %s", tc.key)
			}

			if want := tc.value; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
		})
	}
}

func TestMapLen(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, string]()

	m.Set("a", "a")
	if want, got := 1, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}

	m.Set("b", "b")
	if want, got := 2, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}

	// Key already exists.
	m.Set("a", "1")
	if want, got := 2, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}
}

func TestMapDelete(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, int]()

	key := "key"

	m.Set(key, 1)
	_, ok := m.Get(key)
	if !ok {
		t.Fatalf("key does not exist: %s", key)
	}

	if want, got := 1, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}

	m.Set("keep", 2)
	if want, got := 2, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}

	m.Delete(key)
	_, ok = m.Get(key)
	if ok {
		t.Fatalf("key was not deleted: %s", key)
	}

	if want, got := 1, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}
}

func TestMapContains(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, string]()

	if ok := m.Contains("a"); ok {
		t.Fatal("want key to be absent from empty map")
	}

	m.Set("a", "a")
	if ok := m.Contains("a"); !ok {
		t.Fatal("want key to exist in map after Set")
	}

	m.Delete("a")
	if ok := m.Contains("a"); ok {
		t.Fatal("want key to be absent after Delete")
	}
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

	if want, got := len(cases), m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}

	m.Clear()

	if want, got := 0, m.Len(); want != got {
		t.Fatalf("want %d, got %d", want, got)
	}

	for _, tc := range cases {
		t.Run(tc.key, func(t *testing.T) {
			_, ok := m.Get(tc.key)
			if ok {
				t.Errorf("key was not deleted on clear: %s", tc.key)
			}
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
	if !ok {
		t.Fatalf("key does not exist: %s", key)
	}

	if want := value; want != got {
		t.Errorf("want value %v, got %v", want, got)
	}
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
