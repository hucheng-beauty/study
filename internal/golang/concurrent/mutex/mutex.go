package main

import (
	"log"
	"sync"
)

// Counter 计数器
type Counter struct {
	mutex sync.Mutex
	c     int
}

func (c *Counter) Increase() {
	c.mutex.Lock()
	c.c++
	c.mutex.Unlock()
}

func NewCounter() *Counter {
	return &Counter{mutex: sync.Mutex{}, c: 0}
}

func main() {
	c := NewCounter()
	log.Printf("before, counter:%d\n", c.c)
	for i := 0; i < 10; i++ {
		c.Increase()
	}
	log.Printf("after, counter:%d\n", c.c)
}
