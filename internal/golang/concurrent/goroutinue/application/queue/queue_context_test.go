package toolkit

import (
    "context"
    "log"
    "sync"
    "testing"
    "time"
)

// TestQueueContextCancellation tests that the queue processes remaining jobs
// even when context is cancelled before Close() is called
func TestQueueContextCancellation(t *testing.T) {
    var processed []int
    var mu sync.Mutex

    ctx, cancel := context.WithCancel(context.Background())

    q := NewQueue(&Options{
        Ctx:          ctx,
        BatchWait:    1 * time.Second,
        BatchSize:    5,
        ChanSize:     100,
        RoutineNum:   2,
        DrainTimeout: 3 * time.Second, // Custom drain timeout
        Func: func(jobs []interface{}) {
            mu.Lock()
            defer mu.Unlock()
            log.Printf("Processing batch of %d jobs\n", len(jobs))
            for _, job := range jobs {
                if num, ok := job.(int); ok {
                    processed = append(processed, num)
                }
            }
        },
    })

    q.Start()

    // Send 20 jobs
    for i := 0; i < 20; i++ {
        q.Chan() <- i
    }

    log.Println("Sent 20 jobs, waiting a bit...")
    time.Sleep(500 * time.Millisecond)

    // Cancel context (simulating application shutdown)
    log.Println("Cancelling context...")
    cancel()

    // Wait a bit to let drain complete (DrainTimeout is 3s)
    time.Sleep(4 * time.Second)

    // Close the queue properly
    log.Println("Closing queue...")
    q.Close()

    mu.Lock()
    defer mu.Unlock()

    log.Printf("Total processed: %d jobs\n", len(processed))
    log.Printf("Processed jobs: %v\n", processed)

    if len(processed) < 20 {
        t.Errorf("Expected at least 20 jobs processed, but got %d", len(processed))
    }
}

// TestQueueCloseBeforeContextCancel tests the normal close flow
func TestQueueCloseBeforeContextCancel(t *testing.T) {
    var processed []int
    var mu sync.Mutex

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    q := NewQueue(&Options{
        Ctx:        ctx,
        BatchWait:  1 * time.Second,
        BatchSize:  5,
        ChanSize:   100,
        RoutineNum: 2,
        Func: func(jobs []interface{}) {
            mu.Lock()
            defer mu.Unlock()
            log.Printf("Processing batch of %d jobs\n", len(jobs))
            for _, job := range jobs {
                if num, ok := job.(int); ok {
                    processed = append(processed, num)
                }
            }
        },
    })

    q.Start()

    // Send 20 jobs
    for i := 0; i < 20; i++ {
        q.Chan() <- i
    }

    log.Println("Sent 20 jobs, closing queue...")
    q.Close() // Close before context cancel

    mu.Lock()
    defer mu.Unlock()

    log.Printf("Total processed: %d jobs\n", len(processed))

    if len(processed) != 20 {
        t.Errorf("Expected 20 jobs processed, but got %d", len(processed))
    }
}
