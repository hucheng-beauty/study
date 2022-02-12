package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	生成器
	服务/任务
	同时等待多个服务:俩种方法
			->A->
				   ->C  ->
			->B->
*/

/*
	任务控制:
		非阻塞等待
		超时机制
		任务中断或退出
		优雅退出
*/

// 生成器
func msgGenerate(name string, done chan struct{}) chan string {
	out := make(chan string)
	go func() {
		i := 0
		for {
			select {
			case <-time.After(time.Duration(rand.Intn(5000)) * time.Millisecond):
				out <- fmt.Sprintf("service: %s; message: %d", name, i)
			case <-done:
				fmt.Println("cleaning up")
				time.Sleep(2 * time.Second)
				fmt.Println("clean done")
				done <- struct{}{}
				return
			}
			i++
		}
	}()
	return out
}

func fanIn_new(chs ...chan string) chan string {
	out := make(chan string)

	// first
	/*for _, ch := range chs {
		go func() {
			for {
				out <- <-ch
			}
		}()
	}*/

	// second
	for _, ch := range chs {
		chCopy := ch
		go func() {
			for {
				out <- <-chCopy
			}
		}()
	}

	// third
	for _, ch := range chs {
		go func(in chan string) {
			for {
				out <- <-in
			}
		}(ch)
	}

	// fourth
	/*for i := 0; i < len(chs); i++ {
		go func(ii int) {
			for {
				out <- <-chs[ii]
			}
		}(i)
	}*/

	return out
}

// 适用 channel 数量未知
func fanIn(c1, c2 chan string) chan string {
	out := make(chan string)
	go func() {
		for {
			out <- <-c1
		}
	}()
	go func() {
		for {
			out <- <-c2
		}
	}()
	return out
}

// 适用 channel 数量已知
func fanInBySelect(c1, c2 chan string) chan string {
	out := make(chan string)
	go func() {
		s := ""
		for {
			select {
			case s = <-c1:
				out <- s
			case s = <-c2:
				out <- s
			}
		}
	}()
	return out
}

// 非阻塞等待
func nonBlockingWait(c chan string) (string, bool) {
	select {
	case m := <-c:
		return m, true
	default: // select 中有 default 则为非阻塞等待
		return "", false
	}
}

// 超时等待
func timeoutWait(c chan string, timeout time.Duration) (string, bool) {
	select {
	case m := <-c:
		return m, true
	case <-time.After(timeout):
		return "", false
	}
}

func main() {
	done := make(chan struct{})
	m1 := msgGenerate("m1", done)
	//m2 := msgGenerate("m2")
	//m := fanIn_new(m1, m2)

	for i := 0; i < 5; i++ {
		fmt.Println(<-m1)
		if m, ok := timeoutWait(m1, time.Second); ok {
			fmt.Println(m)
		} else {
			fmt.Println("timeout")
		}
	}
	done <- struct{}{}
	<-done
}
