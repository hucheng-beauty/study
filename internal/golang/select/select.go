package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*
	select 进行调度:
		select 的使用(default 非阻塞式从 channel 获取数据)
		定时器的使用
		在 select 中使用 nil channel
*/

/*
	传统的同步模型
		WaitGroup
		Mutex
		Cond
*/
func generator() chan int {
	out := make(chan int)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			out <- i
			i++
		}
	}()
	return out
}

func worker(id int, c chan int) {
	for n := range c {
		time.Sleep(time.Second)
		fmt.Printf("worker %d received %d\n", id, n)
	}
}

func create_worker(id int) chan<- int {
	c := make(chan int)
	go worker(id, c)
	return c
}

func main() {
	var c1, c2 = generator(), generator()
	var worker = create_worker(0)

	tm := time.After(10 * time.Second)
	tt := time.Tick(time.Second)

	var values []int
	for {
		// 利用 channel nil 的特性,select 阻塞
		var activeWorker chan<- int
		var activeValue int
		if len(values) > 0 {
			activeWorker = worker
			activeValue = values[0]
		}

		select {
		case n := <-c1: // 从 c1 接收 data
			values = append(values, n)
		case n := <-c2: // 从 c2 接受 data
			values = append(values, n)
		case activeWorker <- activeValue: // 消费 data
			values = values[1:]
		case <-time.After(800 * time.Millisecond): // 超时
			fmt.Println("timeout")
		case <-tt: // 定时器的使用
			fmt.Println("len(values):", len(values))
		case <-tm: // 计时器(总的服务时长)
			fmt.Println("good bye")
			return
		}
	}
}
