package main

import (
    "context"
    "fmt"
    "runtime"
    "time"

    "study/internal/golang/concurrent/goroutinue/application/queue"
)

func main() {
    fmt.Println("=== Queue Improvements Demo ===\n")

    ctx, cancel := context.WithCancel(context.Background())

    processed := 0
    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 5,
        QueueLength:  20,
        Func: func(job interface{}) {
            processed++
            fmt.Printf("[%d] Processing: %v\n", processed, job)
            time.Sleep(time.Millisecond * 50)
        },
    })

    fmt.Println("1. Starting queue...")
    q.Start()
    fmt.Printf("   Goroutines: %d\n\n", runtime.NumGoroutine())

    fmt.Println("2. Submitting 10 tasks...")
    for i := 1; i <= 10; i++ {
        q.Chan() <- fmt.Sprintf("Task-%d", i)
    }

    time.Sleep(time.Millisecond * 800)
    fmt.Printf("\n3. Processed %d tasks so far\n", processed)
    fmt.Printf("   Goroutines: %d\n\n", runtime.NumGoroutine())

    fmt.Println("4. Stopping queue (draining remaining tasks)...")
    q.Stop() // 通知调度器停止并处理剩余任务
    cancel() // 通知所有 workers 退出
    q.Wait() // 等待所有 workers 完全退出

    fmt.Printf("\n5. Final status:\n")
    fmt.Printf("   Total processed: %d tasks\n", processed)
    fmt.Printf("   Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Println("\n✓ Queue stopped gracefully!")
}
