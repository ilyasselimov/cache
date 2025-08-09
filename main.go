package main

import "sync"

// Cache — один шард кеша
type Cache struct {
	mu    sync.RWMutex
	cache map[int]int
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[int]int),
	}
}

func (c *Cache) Insert(k, v int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[k] = v
}

func (c *Cache) Remove(k int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, k)
}

func (c *Cache) Lookup(k int) (int, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.cache[k]
	return v, ok
}

// SharedCache — кеш с шардами
type SharedCache struct {
	caches      []*Cache
	countShards int
}

// NewSharedCache — конструктор SharedCache
func NewSharedCache(countShards int) *SharedCache {
	caches := make([]*Cache, countShards)
	for i := 0; i < countShards; i++ {
		caches[i] = NewCache()
	}
	return &SharedCache{
		caches:      caches,
		countShards: countShards,
	}
}

// shardForKey — вычисляем индекс шарда по ключу
func (s *SharedCache) shardForKey(k int) int {
	return k % s.countShards
}

func (s *SharedCache) Insert(k int, v int) {
	shard := s.shardForKey(k)
	s.caches[shard].Insert(k, v)
}

func (s *SharedCache) Remove(k int) {
	shard := s.shardForKey(k)
	s.caches[shard].Remove(k)
}

func (s *SharedCache) Lookup(k int) (int, bool) {
	shard := s.shardForKey(k)
	return s.caches[shard].Lookup(k)
}
