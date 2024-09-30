package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

/*
   动态个数协程并发处理任务
       针对固定个数协程的缺点,协程数量最好能够根据需要处理任务的多少动态变更,
       指定一个并发上限,任务多时增加协程数量,任务少时减少协程数量;

	两种实现:
       自定义令牌池实现
       令牌池维持最大允许并发任务数个令牌,每个任务启动时请求令牌,运行完成返回令牌。
*/

// runDynamicTask
// 最大同时运行 maxTaskNum 个任务处理数据
// 自定义令牌池维持 maxTaskNum 个令牌供竞争
func runDynamicTask(dataSource <-chan int, maxTaskNum int) {
	var wg sync.WaitGroup

	// 初始化令牌池并生产令牌
	tokenPool := make(chan struct{}, maxTaskNum)
	for i := 0; i < maxTaskNum; i++ {
		tokenPool <- struct{}{}
	}

	for data := range dataSource {
		// 先消耗令牌,如果被消费完则阻塞等待其它任务返还令牌
		<-tokenPool

		wg.Add(1)
		go func(data int) {
			// 任务运行完成,返还令牌
			defer func() {
				// do something
				time.Sleep(3 * time.Second)
				tokenPool <- struct{}{}
				wg.Add(-1)
			}()
		}(data)
	}

	wg.Wait()
}

// runSemaphoreTask
// 最大同时运行 maxTaskNum 个任务处理数据
// 使用信号量维持 maxTaskNum 个信号
func runSemaphoreTask(dataSource <-chan int, maxTaskNum int64) {
	var wg sync.WaitGroup
	w := semaphore.NewWeighted(maxTaskNum)

	for data := range dataSource {
		// 先获取信号量,如果被消费完则阻塞等待信号量返还
		_ = w.Acquire(context.TODO(), 1)

		wg.Add(1)
		go func(data int) {
			defer func() {
				// 运行完成返还信号量
				defer w.Release(1)
				wg.Done()
			}()

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}

/*
   unBoundedTask + NumGoroutineMonitor
*/

// runNumGoroutineMonitor 协程数量监控
func runNumGoroutineMonitor() {
	log.Printf("初始协程数量->%d\n", runtime.NumGoroutine())

	for {
		select {
		case <-time.After(time.Second):
			log.Printf("协程数量->%d\n", runtime.NumGoroutine())
		}
	}
}

// runTaskDataGenerator 产生数据
func runTaskDataGenerator(dataSource chan int) {
	for i := 0; i < 100; i++ {
		dataSource <- i
	}

	close(dataSource)
}

// runInfiniteTask 每来一个数据起协程处理任务
func runInfiniteTask(dataSource <-chan int) {
	var wg sync.WaitGroup

	for data := range dataSource {
		wg.Add(1)

		go func(data int) {
			defer wg.Add(-1)

			// do something
			time.Sleep(100 * time.Millisecond)
		}(data)
	}

	wg.Wait()
}

func Run() {
	dataSource := make(chan int)

	go runTaskDataGenerator(dataSource)
	go runNumGoroutineMonitor()

	runInfiniteTask(dataSource)
}
