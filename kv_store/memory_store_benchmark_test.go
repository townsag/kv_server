package kv_store

import (
	"testing"
    "fmt"
)

func BenchmarkSet(b *testing.B) {
	var store Store = NewMemoryStore()
	i := 0

	for b.Loop() {
		key := fmt.Sprintf("key-%d", i)
        value := fmt.Sprintf("value-%d", i)
        store.Set(key, value)
	}
}

func BenchmarkSetParallel(b *testing.B) {
    var store = NewMemoryStore()
    b.RunParallel(func(p *testing.PB) {
        i := 0
        for p.Next() {
            key := fmt.Sprintf("key-%d", i)
            value := fmt.Sprintf("value-%d", i)
            store.Set(key, value)
        }
    })
}

func BenchmarkGet(b *testing.B) {
    var store Store = NewMemoryStore()
    store.Set("asdf", "asdf")

    for b.Loop() {
        store.Get("asdf")
    }
}

func BenchmarkDelete(b *testing.B) {
	// prep populate the store with initial data
    var store Store = NewMemoryStore()
    for i := 0; i < 100_000; i++ {
        key := fmt.Sprintf("key-%d", i)
        value := fmt.Sprintf("value-%d", i)

        store.Set(key, value)
    }

    i := 0
    for b.Loop() {
        key := fmt.Sprintf("key-%d", i)
        store.Delete(key)
    }
}