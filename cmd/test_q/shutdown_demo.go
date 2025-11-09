package main

import (
    "context"
    "fmt"
    "time"

    "study/internal/golang/concurrent/goroutinue/application/queue"
)

func main() {
    fmt.Println("=== 演示不同的关闭方式 ===\n")

    // 方式 1: 手动三步关闭（灵活但繁琐）
    demo1ManualShutdown()

    // 方式 2: 使用 Shutdown 一站式关闭（推荐）
    demo2ShutdownMethod()
}

func demo1ManualShutdown() {
    fmt.Println("【方式1】手动三步关闭")
    ctx, cancel := context.WithCancel(context.Background())

    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 3,
        QueueLength:  10,
        Func: func(job interface{}) {
            fmt.Printf("  处理: %v\n", job)
            time.Sleep(time.Millisecond * 50)
        },
    })

    q.Start()

    // 提交任务
    for i := 1; i <= 5; i++ {
        q.Chan() <- fmt.Sprintf("Task-%d", i)
    }

    time.Sleep(time.Millisecond * 200)

    // 手动三步关闭
    fmt.Println("  步骤1: 调用 Stop()...")
    q.Stop()

    fmt.Println("  步骤2: 调用 cancel()...")
    cancel()

    fmt.Println("  步骤3: 调用 Wait()...")
    q.Wait()

    fmt.Println("  ✓ 关闭完成\n")
}

func demo2ShutdownMethod() {
    fmt.Println("【方式2】使用 Shutdown 一站式关闭（推荐）")
    ctx, cancel := context.WithCancel(context.Background())

    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 3,
        QueueLength:  10,
        Func: func(job interface{}) {
            fmt.Printf("  处理: %v\n", job)
            time.Sleep(time.Millisecond * 50)
        },
    })

    q.Start()

    // 提交任务
    for i := 1; i <= 5; i++ {
        q.Chan() <- fmt.Sprintf("Task-%d", i)
    }

    time.Sleep(time.Millisecond * 200)

    // 一站式关闭！
    fmt.Println("  调用 Shutdown(cancel)...")
    q.Shutdown(cancel)

    fmt.Println("  ✓ 关闭完成\n")
}
