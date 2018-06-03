package hn

import (
	"fmt"
	"sync"
	"testing"
)

func TestCache(t *testing.T) {
	cache := &Cache{
		lock:    &sync.RWMutex{},
		stories: make(map[int]*Item),
	}
	cache.Get(3)
}

func ExampleCache() {
	newCache()
	id, err := cache.Get(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}

// Example of get from cache
func ExampleCache_Get() *Item {
	id := 3
	item, err := cache.Get(id)
	if err != nil {
		panic(err)
	}
	return item
}

func ExampleCache_Put() {
	item := &Item{}
	id := 1
	cache.Put(id, item)
}
