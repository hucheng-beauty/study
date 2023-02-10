package main

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

/*
	指定处理速度并发处理任务
        针对固定个数协程的缺点,借鉴限流器的实现,控制每个时刻最大允许协程数量也达到控制协程数量的目的

	俩种实现:
        自定义令牌池实现
        相当于一个简单限流器,指定速度生产令牌,每个任务启动时必须请求到令牌;
*/

// runCustomTokenPoolTask 限制每秒允许的最大协程数量,限流器的思路
func runCustomTokenPoolTask(dataSource <-chan int) {
	var wg sync.WaitGroup

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

	for data := range dataSource {
		// 先获取令牌,如果被消费完则阻塞等待新令牌产生
		<-tokenPool

		wg.Add(1)
		go func(data int) {
			defer wg.Add(-1)

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}

// runRateLimitTask 限制每秒允许的最大协程数量,使用官方限流器
func runRateLimitTask(dataSource <-chan int) {
	var wg sync.WaitGroup

	// 初始化令牌池
	limiter := rate.NewLimiter(rate.Every(time.Second) /*每秒一个*/, 10)

	for data := range dataSource {
		// 先获取令牌,如果被消费完则阻塞等待新令牌产生
		_ = limiter.Wait(context.TODO())

		wg.Add(1)
		go func(data int) {
			defer wg.Add(-1)

			// do something
			time.Sleep(3 * time.Second)
		}(data)
	}

	wg.Wait()
}
