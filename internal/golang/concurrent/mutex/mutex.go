package main

import (
	"log"
	"sync"
	"time"
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
		go c.Increase()
	}

	time.Sleep(100 * time.Millisecond)
	log.Printf("after, counter:%d\n", c.c)

	// data race
	number := 0
	for i := 0; i < 10; i++ {
		go func() {
			number++
		}()
	}

}
