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

func createWorkerSecond(id int, c chan int) {
	// channel range
	for n := range c {
		fmt.Printf("worker %d received %c\n", id, n)
	}
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

func nonBlockingWaiting(c chan string) (string, bool) {
	select {
	case m := <-c:
		return m, true
	// select 中有 default 则为非阻塞等待
	default:
		return "", false
	}
}

func timeoutWaiting(c chan string, timeout time.Duration) (string, bool) {
	select {
	case m := <-c:
		return m, true
	case <-time.After(timeout):
		return "", false
	}
}

func PrintNumberAndChar() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{}, 1)

	ch2 <- struct{}{}
	go func() {
		for i := 1; i <= 26; i++ {
			<-ch1
			fmt.Printf("%d\n", i)
			ch1 <- struct{}{}
		}
	}()

	go func() {
		for i := 'a'; i <= 'z'; i++ {
			<-ch1
			fmt.Printf("%s\n", string(i))
			ch2 <- struct{}{}
		}
	}()

	time.Sleep(100 * time.Millisecond)
}

func Print(id int, ch chan struct{}, nextChan chan struct{}) {
	if id%2 == 0 {
		for {
			for i := 1; i <= 26; i++ {
				<-ch
				fmt.Printf("%d\n", i)
				time.Sleep(time.Second)
				nextChan <- struct{}{}
			}
		}
	} else {
		for {
			for i := 'a'; i <= 'z'; i++ {
				<-ch
				fmt.Printf("%s\n", string(i))
				nextChan <- struct{}{}
			}
		}
	}
}

func PrintNumberAndCharSecond() {
	chs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
	}

	for i := 0; i < 2; i++ {
		go Print(i, chs[i], chs[(i+1)%2])
	}

	chs[0] <- struct{}{}
	time.Sleep(10 * time.Second)
}

func main() {
	// unBuffered and buffered channel

	// unBufferedChan()
	// bufferedChan()

	// closed channel
	//closedChan()

	PrintNumberAndCharSecond()
}
