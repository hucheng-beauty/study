package main

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/*
	指定处理速度并发处理任务
        针对固定个数协程的缺点，另一个思路是借鉴限流器的实现，
        控制每个时刻最大允许协程数量也达到控制协程数量的目的。这里也提供两种实现思路

        自定义令牌池实现
        相当于一个简单限流器，指定速度生产令牌，每个任务启动时必须请求到令牌。
*/

// runRateLimitTask 限制每秒允许的最大协程数量，限流器的思路
func runRateLimitTask(dataChan <-chan int) {
	// 初始化令牌池
	tokenPool := make(chan struct{})
	go func() {
		for {
			select {
			// 动态控制令牌生成速度
			case <-time.After(time.Second):
				tokenPool <- struct{}{}
			}
		}
	}()

	var wg sync.WaitGroup

	for data := range dataChan {
		// 先获取令牌，如果被消费完则阻塞等待新令牌产生
		<-tokenPool

		wg.Add(1)
		go func(data int) {
			defer wg.Done()

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}

// runRateLimitTask2 限制每秒允许的最大协程数量，使用官方限流器
func runRateLimitTask2(dataChan <-chan int) {
	// 初始化令牌池
	limit := rate.Every(time.Second) // 每秒一个
	limiter := rate.NewLimiter(limit, 10)

	var wg sync.WaitGroup

	for data := range dataChan {
		// 先获取令牌，如果被消费完则阻塞等待新令牌产生
		_ = limiter.Wait(context.TODO())

		wg.Add(1)
		go func(data int) {
			defer wg.Done()

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}
