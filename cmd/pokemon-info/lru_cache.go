package main

import "container/list"

// lruCap bounds the in-memory session cache. Entries are small API
// structs; eviction only kicks in past the cap, so short sessions never
// pay any eviction cost.
const lruCap = 100

type lruEntry struct {
	key string
	val any
}

// lruCache is a fixed-capacity map with least-recently-used eviction.
type lruCache struct {
	cap int
	m   map[string]*list.Element
	ll  *list.List
}

func newLRUCache(cap int) *lruCache {
	return &lruCache{cap: cap, m: map[string]*list.Element{}, ll: list.New()}
}

func (c *lruCache) get(key string) (any, bool) {
	if e, ok := c.m[key]; ok {
		c.ll.MoveToFront(e)
		return e.Value.(lruEntry).val, true
	}
	return nil, false
}

func (c *lruCache) has(key string) bool {
	_, ok := c.m[key]
	return ok
}

func (c *lruCache) set(key string, val any) {
	if e, ok := c.m[key]; ok {
		e.Value = lruEntry{key: key, val: val}
		c.ll.MoveToFront(e)
		return
	}
	e := c.ll.PushFront(lruEntry{key: key, val: val})
	c.m[key] = e
	for c.ll.Len() > c.cap {
		oldest := c.ll.Back()
		if oldest == nil {
			break
		}
		c.ll.Remove(oldest)
		delete(c.m, oldest.Value.(lruEntry).key)
	}
}

func (c *lruCache) clear() {
	c.m = map[string]*list.Element{}
	c.ll = list.New()
}
