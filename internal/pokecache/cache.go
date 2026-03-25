package pokecache

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	m        map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
}

func NewCache(i time.Duration) *Cache {
	newCache := &Cache{
		m:        make(map[string]cacheEntry),
		interval: i,
	}
	go newCache.reapLoop()
	return newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = cacheEntry{
		time.Now(),
		val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if res, ok := c.m[key]; ok == true {
		return res.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.m {
			if entry.createdAt.Add((c.interval)).Before(time.Now()) {
				delete(c.m, key)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache) Check(url string) ([]byte, error) {
	var body []byte
	if val, ok := c.Get(url); ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return body, err
		}
		defer res.Body.Close()
		if res.StatusCode >= 300 {
			return body, fmt.Errorf("Error: %d, check spelling", res.StatusCode)
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return body, err
		}
		c.Add(url, body)
	}
	return body, nil
}
