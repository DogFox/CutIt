package cache

import (
	"container/list"
	"sync"
)

type Cache struct {
	maxSize int
	items   map[string]*list.Element
	order   *list.List
	mu      sync.Mutex
}

type entry struct {
	key   string
	value string
}

func NewCache(size int) *Cache {
	return &Cache{
		maxSize: size,
		items:   make(map[string]*list.Element),
		order:   list.New(),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, found := c.items[key]; found {
		c.order.MoveToFront(el)
		return el.Value.(*entry).value, true
	}
	return "", false
}

func (c *Cache) Put(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, found := c.items[key]; found {
		c.order.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}

	if c.order.Len() >= c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			delete(c.items, oldest.Value.(*entry).key)
			c.order.Remove(oldest)
		}
	}

	el := c.order.PushFront(&entry{key, value})
	c.items[key] = el
}
