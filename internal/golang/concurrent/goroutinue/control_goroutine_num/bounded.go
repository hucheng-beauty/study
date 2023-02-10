package main

import (
	"sync"
	"time"
)

/*
   固定个数协程并发处理任务
       一般叫做Bounded/Fixed并发控制。

       优点是简单，不复杂的并发任务这样简单处理即可。
       缺点在于dataChan可能流量不不均衡，需要同时处理的任务多少在变动，
       但是对应的协程数量保持不变，要不就是任务处理堵塞要不就是存在多余的协程空闲。
*/

// runBoundedTask 起 maxTaskNum 个协程共同处理任务
func runBoundedTask(dataChan <-chan int, maxTaskNum int) {
	var wg sync.WaitGroup
	wg.Add(maxTaskNum)

	for i := 0; i < maxTaskNum; i++ {
		go func() {
			defer wg.Done()

			for data := range dataChan {
				func(data int) {

					// do something
					time.Sleep(3 * time.Second)
				}(data)
			}
		}()
	}

	wg.Wait()
}
