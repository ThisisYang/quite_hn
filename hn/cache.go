package hn

import (
	"errors"
	"sync"
)

var cache *Cache

func init() {
	newCache()
}

type Cache struct {
	lock    *sync.RWMutex
	stories map[int]*Item
}

func newCache() {
	cache = &Cache{
		lock:    &sync.RWMutex{},
		stories: make(map[int]*Item),
	}
}

func (c *Cache) Put(id int, i *Item) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.stories[id] = i
}

func (c *Cache) Get(id int) (*Item, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	i, ok := c.stories[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return i, nil
}
