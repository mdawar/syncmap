package syncmap_test

import (
	"fmt"

	"github.com/mdawar/syncmap"
)

func Example() {
	m := syncmap.New[string, int]()

	m.Set("a", 1)
	m.Set("b", 1)

	m.Delete("b")

	fmt.Println(m.Len())
	fmt.Println(m.Get("a"))
	fmt.Println(m.Get("b"))

	m.Clear()

	fmt.Println(m.Len())

	// Output:
	// 1
	// 1 true
	// 0 false
	// 0
}

func ExampleNewWithCapacity() {
	// Create a map with a capacity hint.
	m := syncmap.NewWithCapacity[string, int](10_000)

	m.Set("a", 1)

	fmt.Println(m.Len())
	fmt.Println(m.Get("a"))
	fmt.Println(m.Get("b"))

	// Output:
	// 1
	// 1 true
	// 0 false
}

func ExampleMap_All() {
	m := syncmap.New[string, int]()

	for k, v := range m.All() {
		fmt.Println("Key:", k, "-", "Value:", v)
	}
}
