package main

import (
    "context"
    "fmt"
    "runtime"
    "time"

    "study/internal/golang/concurrent/goroutinue/application/queue"
)

func main() {
    fmt.Println("=== 测试改进后的队列 ===\n")

    // 测试 1: 基本功能
    test1BasicFunction()

    // 测试 2: 优雅退出
    test2GracefulShutdown()

    // 测试 3: Context Done 响应
    test3ContextResponse()
}

func test1BasicFunction() {
    fmt.Println("【测试1】基本功能测试")
    ctx, cancel := context.WithCancel(context.Background())

    processed := 0
    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 3,
        QueueLength:  10,
        Func: func(job interface{}) {
            processed++
            fmt.Printf("  处理任务: %v\n", job)
            time.Sleep(time.Millisecond * 50)
        },
    })

    fmt.Println("  启动队列...")
    q.Start()

    fmt.Println("  提交5个任务...")
    for i := 1; i <= 5; i++ {
        q.Chan() <- fmt.Sprintf("Task-%d", i)
    }

    time.Sleep(time.Millisecond * 500)
    fmt.Printf("  已处理 %d 个任务\n", processed)

    // 使用 Shutdown 一站式关闭
    fmt.Println("  使用 Shutdown 关闭...")
    q.Shutdown(cancel)

    fmt.Printf("  最终处理: %d 个任务\n", processed)
    fmt.Printf("  Goroutines: %d\n\n", runtime.NumGoroutine())
}

func test2GracefulShutdown() {
    fmt.Println("【测试2】优雅退出测试（处理剩余任务）")
    ctx, cancel := context.WithCancel(context.Background())

    processed := 0
    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 2,
        QueueLength:  20,
        Func: func(job interface{}) {
            processed++
            fmt.Printf("  处理任务: %v\n", job)
            time.Sleep(time.Millisecond * 100)
        },
    })

    q.Start()

    fmt.Println("  提交10个任务...")
    for i := 1; i <= 10; i++ {
        q.Chan() <- fmt.Sprintf("Task-%d", i)
    }

    time.Sleep(time.Millisecond * 150) // 只处理一部分
    fmt.Printf("  部分处理: %d 个任务\n", processed)

    fmt.Println("  使用 Shutdown 优雅关闭...")
    q.Shutdown(cancel)

    fmt.Printf("  最终处理: %d 个任务\n", processed)
    fmt.Printf("  Goroutines: %d\n\n", runtime.NumGoroutine())
}

func test3ContextResponse() {
    fmt.Println("【测试3】Context Done 响应测试")
    ctx, cancel := context.WithCancel(context.Background())

    processed := 0
    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 3,
        QueueLength:  10,
        Func: func(job interface{}) {
            processed++
            fmt.Printf("  处理任务: %v\n", job)
            time.Sleep(time.Millisecond * 50)
        },
    })

    q.Start()

    fmt.Println("  提交3个任务...")
    for i := 1; i <= 3; i++ {
        q.Chan() <- fmt.Sprintf("Task-%d", i)
    }

    time.Sleep(time.Millisecond * 200)

    fmt.Println("  直接取消 context...")
    cancel() // 直接取消，不调用 Stop()
    q.Wait()

    fmt.Printf("  最终处理: %d 个任务\n", processed)
    fmt.Printf("  Goroutines: %d\n\n", runtime.NumGoroutine())
}
