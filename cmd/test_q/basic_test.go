package main

import (
    "context"
    "fmt"
    "time"

    "study/internal/golang/concurrent/goroutinue/application/queue"
)

func main() {
    fmt.Println("Starting basic test...")

    ctx, cancel := context.WithCancel(context.Background())

    q := queue.NewQueue(&queue.Options{
        Ctx:          ctx,
        WorkersCount: 2,
        QueueLength:  5,
        Func: func(job interface{}) {
            fmt.Printf("Job: %v\n", job)
        },
    })

    fmt.Println("Queue created")
    q.Start()
    fmt.Println("Queue started")

    q.Chan() <- "test1"
    fmt.Println("Sent test1")

    time.Sleep(time.Millisecond * 100)
    fmt.Println("Slept 100ms")

    fmt.Println("Calling Stop()...")
    q.Stop()
    fmt.Println("Stop() called")

    fmt.Println("Calling cancel()...")
    cancel()
    fmt.Println("cancel() called")

    fmt.Println("Calling Wait()...")
    q.Wait()
    fmt.Println("Wait() returned")

    fmt.Println("DONE")
}
