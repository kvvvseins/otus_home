package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mu       *sync.Mutex
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mu:       &sync.Mutex{},
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val, ok := c.items[key]; ok {
		val.Value = value
		c.queue.MoveToFront(val)

		return true
	}

	if c.queue.Len() >= c.capacity {
		lastItem := c.queue.Back()

		c.queue.Remove(lastItem)

		for k, v := range c.items {
			if v == lastItem {
				delete(c.items, k)

				break
			}
		}
	}

	c.queue.PushFront(value)
	c.items[key] = c.queue.Front()

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if value, ok := c.items[key]; ok {
		c.queue.MoveToFront(value)

		return value.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
