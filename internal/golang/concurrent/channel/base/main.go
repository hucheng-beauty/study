package main

import (
	"fmt"
	"time"
)

func closedChan() {
	ch := make(chan int)
	go func() {
		for {
			// receive
			value, ok := <-ch
			if !ok {
				return
			}
			fmt.Printf("received value: %d\n", value)
		}
	}()

	// send
	for i := 0; i < 3; i++ {
		ch <- 'a' + i
	}

	// 发送方主动关闭 channel
	time.Sleep(time.Millisecond)
	close(ch)
}

func bufferedChan() {
	bufferChan := make(chan int, 3)
	// send
	for i := 0; i < 3; i++ {
		bufferChan <- 'a' + i
	}

	// receive
	go func() {
		for {
			value, ok := <-bufferChan
			if !ok {
				return
			}
			fmt.Printf("received value: %d\n", value)
		}
	}()

	time.Sleep(time.Millisecond)
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

func unBufferedChan() {
	// 接受方开启死循环一直接收,主协程发送完主协程结束,接受方都退出
	var channels [10]chan<- int

	// receive
	for i := 0; i < 10; i++ {
		channels[i] = createWorker(i)
	}

	// send
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}

	time.Sleep(time.Millisecond)
}

func main() {
	// unBuffered and buffered channel

	unBufferedChan()
	bufferedChan()

	// closed channel
	// closedChan()

}
