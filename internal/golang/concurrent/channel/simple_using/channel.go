package main

import (
	"fmt"
	"time"
)

func closedChan() {
	ch := make(chan int)
	go func() {
		for {
			// read
			value, ok := <-ch
			if !ok {
				return
			}
			fmt.Printf("received value: %d\n", value)
		}
	}()

	// write
	for i := 0; i < 3; i++ {
		ch <- 'a' + i
	}

	// 发送方 close
	time.Sleep(time.Millisecond)
	close(ch)
}

func bufferedChan() {
	bufferChan := make(chan int, 3)
	// write
	for i := 0; i < 3; i++ {
		bufferChan <- 'a' + i
	}

	go func() {
		for {
			// read
			value, ok := <-bufferChan
			if !ok {
				return
			}
			fmt.Printf("received value: %d\n", value)
		}
	}()

	time.Sleep(time.Millisecond)
}

func workerSecond(id int, c chan int) {
	// channel range
	for n := range c {
		fmt.Printf("worker %d received %c\n",
			id, n)
	}
}

func worker(id int, c chan int) {

}

func createWorker(id int) chan<- int {
	c := make(chan int)
	go func() {
		for {
			n, ok := <-c
			if !ok {
				return
			}
			fmt.Printf("worker %d received %c\n", id, n)
		}
	}()
	return c
}

func demoChan() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = createWorker(i)
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}

	time.Sleep(time.Millisecond)
}

func main() {
	// demo channel
	//demoChan()

	// buffered channel
	// bufferedChan()

	// closed channel
	closedChan()
}
