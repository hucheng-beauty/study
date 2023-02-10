package main

import (
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// runGoroutinePoolTask 使用协程池动态管理协程数量
func runGoroutinePoolTask(dataChan <-chan int, maxTaskNum int) {
	var p, _ = ants.NewPool(maxTaskNum)
	defer p.Release()

	var wg sync.WaitGroup

	for _ = range dataChan {
		wg.Add(1)

		// 提交任务，协程池动态管理数量，可以做更多的分配优化策略
		_ = p.Submit(func() {
			defer wg.Done()

			// do something
			time.Sleep(3 * time.Second)
		})

	}

	wg.Wait()
}
