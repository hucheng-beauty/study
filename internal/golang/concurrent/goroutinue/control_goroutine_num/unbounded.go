package main

import (
	"log"
	"runtime"
	"sync"
	"time"
)

// runNumGoroutineMonitor 协程数量监控
func runNumGoroutineMonitor() {
	log.Printf("协程数量->%d\n", runtime.NumGoroutine())

	for {
		select {
		case <-time.After(time.Second):
			log.Printf("协程数量->%d\n", runtime.NumGoroutine())
		}
	}
}

// runTaskDataGenerator 产生数据
func runTaskDataGenerator(dataChan chan int) {
	for i := 0; i < 100; i++ {
		dataChan <- i
	}

	close(dataChan)
}

// runInfiniteTask 每来一个数据起协程处理任务
func runInfiniteTask(dataChan <-chan int) {
	var wg sync.WaitGroup

	for data := range dataChan {
		wg.Add(1)
		go func(data int) {
			defer wg.Done()

			// do something
			time.Sleep(100 * time.Millisecond)
		}(data)
	}

	wg.Wait()
}

func Run() {
	dataChan := make(chan int)

	go runTaskDataGenerator(dataChan)
	go runNumGoroutineMonitor()

	runInfiniteTask(dataChan)
}
