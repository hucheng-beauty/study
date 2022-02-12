package main

import (
	"fmt"
	"time"
)

/*
	channel:CSP(Communication Sequential Process)
*/

func worker_second(id int, c chan int) {
	// channel range
	for n := range c {
		fmt.Printf("worker %d received %c\n",
			id, n)
	}
}

func worker(id int, c chan int) {
	for {
		n, ok := <-c
		if !ok {
			return
		}
		fmt.Printf("worker %d received %c\n",
			id, n)
	}
}

func create_worker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func chanDemo() {
	var channels [10]chan<- int
	for i := 0; i < 10; i++ {
		channels[i] = create_worker(i)
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}

	time.Sleep(time.Millisecond)
}

func bufferedChannel() {
	c := make(chan int, 3)
	c <- 'a'
	c <- 'b'
	c <- 'c'

	go worker(0, c)
	time.Sleep(time.Millisecond)
}

func channelClosed() {
	c := make(chan int)
	go worker(0, c)

	c <- 'a'
	c <- 'b'
	c <- 'c'
	close(c) // 发送方 close
	time.Sleep(time.Millisecond)
}

func main() {
	fmt.Println("channel as first citizen")
	chanDemo()

	fmt.Println("buffered channel")
	bufferedChannel()

	fmt.Println("channel closed and range channel")
	channelClosed()
}
