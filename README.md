# syncmap

[![Go Reference](https://pkg.go.dev/badge/github.com/mdawar/syncmap.svg)](https://pkg.go.dev/github.com/mdawar/syncmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/mdawar/syncmap)](https://goreportcard.com/report/github.com/mdawar/syncmap)
[![Tests](https://github.com/mdawar/syncmap/actions/workflows/test.yml/badge.svg)](https://github.com/mdawar/syncmap/actions/workflows/test.yml)

A simple and generic **Go** map that is safe for concurrent use.

## Installation

```sh
go get -u github.com/mdawar/syncmap
```

## Usage

```go
// Create a map that is safe for concurrent use.
m := syncmap.New[string, int]()

// Create.
m.Set("a", 1)
m.Set("b", 2)

// Get stored value.
v, ok := m.Get("a")
fmt.Println(v)  // 1
fmt.Println(ok) // true

// Delete.
m.Delete("b")

// Map length.
m.Len()

// Clear the map (Remove all entries).
m.Clear()

// Iteration.
for k, v := range m.All() {
  fmt.Println("Key", k, "/", "Value", v)
}
```

```go
// Create a map with an initial capacity hint.
m := syncmap.NewWithCapacity[string, int](10_000)

// Equivalent to.
make(map[string]int, 10_000)
```

## Tests

```sh
go test -race -cover -vet=all
# If you have "just" installed.
just test
# Or using make.
make test
```
