package main

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/semaphore"
)

/*
   动态个数协程并发处理任务
       针对固定个数协程的缺点，一个思路是协程数量最好能够根据来的处理任务的多少，动态变更，
       指定一个并发上限，任务多时增加协程数量，任务少时减少协程数量。这里提供两种实现思路

       自定义令牌池实现
       令牌池维持最大允许并发任务数个令牌，每个任务启动时请求令牌，运行完成返回令牌。
*/

// runDynamicTask
// 最大同时运行 maxTaskNum 个任务处理数据
// 自定义令牌池维持 maxTaskNum 个令牌供竞争
func runDynamicTask(dataChan <-chan int, maxTaskNum int) {
	// 初始化令牌池
	tokenPool := make(chan struct{}, maxTaskNum)
	// 生产令牌
	for i := 0; i < maxTaskNum; i++ {
		tokenPool <- struct{}{}
	}

	var wg sync.WaitGroup

	for data := range dataChan {
		// 先获取令牌，如果被消费完则阻塞等待其它任务返还令牌
		<-tokenPool

		wg.Add(1)
		go func(data int) {
			// 任务运行完成，返还令牌
			defer func() {
				tokenPool <- struct{}{}
				wg.Done()
			}()

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}

// runSemaphoreTask
// 最大同时运行 maxTaskNum 个任务处理数据
// 使用信号量维持 maxTaskNum 个信号
func runSemaphoreTask(dataChan <-chan int, maxTaskNum int64) {
	w := semaphore.NewWeighted(maxTaskNum)

	var wg sync.WaitGroup

	for data := range dataChan {
		// 先获取信号量，如果被消费完则阻塞等待信号量返还
		_ = w.Acquire(context.TODO(), 1)

		wg.Add(1)
		go func(data int) {
			defer wg.Done()

			// 运行完成返还信号量
			defer w.Release(1)

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}
