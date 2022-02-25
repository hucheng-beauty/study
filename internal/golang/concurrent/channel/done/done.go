package main

import (
	"fmt"
)

type worker struct {
	in   chan int
	done chan bool
}

func doWorker(id int, c chan int, done chan bool) {
	for n := range c {
		fmt.Printf("worker %d received %c\n", id, n)
		done <- true
		//go func() { done <- true }()
	}
}

func createWorker(id int) worker {
	w := worker{in: make(chan int), done: make(chan bool)}
	go doWorker(id, w.in, w.done)
	return w
}

func chanDemo() {
	var workers [10]worker
	// crete_worker
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i)
	}

	// write 'a'
	for i, worker := range workers {
		worker.in <- 'a' + i
	}

	// write 'A'
	for i, worker := range workers {
		worker.in <- 'A' + i
	}

	// wait all channels
	for _, worker := range workers {
		<-worker.done
		<-worker.done
	}
}

func main() {
	chanDemo()
}
