package main

import (
	"fmt"
	"sync"
	"time"
)

type AtomicInt struct {
	value int
	lock  sync.Mutex
}

func (a *AtomicInt) increase() {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.value++
}

func (a *AtomicInt) get() int {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.value
}

func main() {
	var ma AtomicInt
	ma.increase()
	go func() {
		ma.increase()
	}()

	time.Sleep(time.Millisecond)
	fmt.Println(ma.get())
}
