package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"golang.org/x/sync/semaphore"
)

func main() {
	maxWorkers := runtime.GOMAXPROCS(0)              // worker 数量
	sema := semaphore.NewWeighted(int64(maxWorkers)) // 信号量
	task := make([]int, maxWorkers*4)                // 任务数，是worker的四倍

	ctx := context.Background()

	for i := range task {
		// 如果没有 worker 可用，会阻塞在这里，直到某个 worker 被释放
		if err := sema.Acquire(ctx, 1); err != nil {
			break
		}

		// 启动 worker goroutine
		go func(i int) {
			defer sema.Release(1)
			time.Sleep(100 * time.Millisecond) // 模拟一个耗时操作
			task[i] = i + 1
		}(i)
	}

	// 请求所有的 worker,这样能确保前面的 worker 都执行完
	if err := sema.Acquire(ctx, int64(maxWorkers)); err != nil {
		log.Printf("获取所有的worker失败: %v", err)
	}

	fmt.Println(task)
}
