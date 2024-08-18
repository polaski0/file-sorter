package main

import "sync"

type Counter struct {
	mu sync.Mutex
	v  map[string]int
}

func NewCounter() *Counter {
	return &Counter{
		v: make(map[string]int),
	}
}

func (c *Counter) Add(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	curr := 0
	if ctr, ok := c.v[key]; ok {
        curr = ctr + 1
		c.v[key] = curr
	} else {
		curr = 1
		c.v[key] = curr
	}
	return curr
}
