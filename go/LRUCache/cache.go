package lrucache

import (
	"container/list"
	"sync"
)

// LRU Cache. See: https://en.wikipedia.org/wiki/Cache_replacement_policies#Least_Recently_Used_(LRU)
type Cache struct {
	lock      sync.Mutex
	capacity  int
	items     map[string]*list.Element
	evictList *list.List
}

type cacheEntry struct {
	key   string
	value any
}

// Creates a new Cache with given capacity
func NewCache(capacity int) *Cache {
	return &Cache{
		capacity:  capacity,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Reports whether an entry with the specified key was found and returns the value if it exists
func (c *Cache) Get(key string) (any, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, ok := c.items[key]; ok {
		c.evictList.MoveToFront(element)
		return element.Value.(*cacheEntry).value, true
	}
	return nil, false
}

// Assign value to key
func (c *Cache) Save(key string, value any) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, ok := c.items[key]; ok {
		element.Value.(*cacheEntry).value = value
		c.evictList.MoveToFront(element)
		return
	}

	entry := &cacheEntry{key: key, value: value}
	element := c.evictList.PushFront(entry)
	c.items[key] = element

	if c.evictList.Len() > c.capacity {
		c.evict()
	}
}

func (c *Cache) evict() {
	element := c.evictList.Back()
	if element != nil {
		c.evictList.Remove(element)
		entry := element.Value.(*cacheEntry)
		delete(c.items, entry.key)
	}
}
