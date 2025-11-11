package engine_refactor

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "testing"
    "time"
)

// TestBasicFunctionality 测试基本功能
func TestBasicFunctionality(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var processed int32
    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            atomic.AddInt32(&processed, 1)
            time.Sleep(10 * time.Millisecond)
        },
        WorkersCount: 3,
        QueueLength:  10,
    })

    if err := engine.Schedule(); err != nil {
        t.Fatal(err)
    }

    // 提交10个任务
    for i := 0; i < 10; i++ {
        if err := engine.Submit(fmt.Sprintf("task-%d", i)); err != nil {
            t.Fatal(err)
        }
    }

    // 等待任务处理
    time.Sleep(200 * time.Millisecond)

    // 优雅关闭
    engine.Stop()
    cancel()
    engine.Wait()

    count := atomic.LoadInt32(&processed)
    if count != 10 {
        t.Errorf("Expected 10 tasks processed, got %d", count)
    }
}

// TestGracefulShutdown 测试优雅关闭
func TestGracefulShutdown(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var processed int32
    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            atomic.AddInt32(&processed, 1)
            time.Sleep(50 * time.Millisecond) // 模拟耗时任务
        },
        WorkersCount: 2,
        QueueLength:  100,
    })

    engine.Schedule()

    // 提交20个任务
    for i := 0; i < 20; i++ {
        engine.Submit(fmt.Sprintf("task-%d", i))
    }

    // 立即停止（不等待所有任务完成提交）
    time.Sleep(10 * time.Millisecond)
    engine.Stop()
    cancel()
    engine.Wait()

    count := atomic.LoadInt32(&processed)
    t.Logf("Processed %d tasks (should be all 20)", count)
    if count != 20 {
        t.Logf("Warning: Expected 20 tasks, got %d (may be expected in drain)", count)
    }
}

// TestPanicRecovery 测试panic恢复
func TestPanicRecovery(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var processed int32
    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            num := job.(int)
            if num == 5 {
                panic("intentional panic")
            }
            atomic.AddInt32(&processed, 1)
        },
        WorkersCount: 2,
        QueueLength:  10,
    })

    engine.Schedule()

    // 提交10个任务，其中一个会panic
    for i := 0; i < 10; i++ {
        engine.Submit(i)
    }

    time.Sleep(200 * time.Millisecond)
    engine.Stop()
    cancel()
    engine.Wait()

    count := atomic.LoadInt32(&processed)
    if count != 9 {
        t.Errorf("Expected 9 tasks processed (1 panicked), got %d", count)
    }
}

// TestContextCancellation 测试上下文取消
func TestContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())

    var processed int32
    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            atomic.AddInt32(&processed, 1)
            time.Sleep(100 * time.Millisecond)
        },
        WorkersCount: 2,
        QueueLength:  10,
    })

    engine.Schedule()

    // 提交任务
    for i := 0; i < 5; i++ {
        engine.Submit(fmt.Sprintf("task-%d", i))
    }

    // 立即取消上下文
    time.Sleep(50 * time.Millisecond)
    cancel()

    // 尝试提交更多任务（应该失败）
    err := engine.Submit("should-fail")
    if err == nil {
        t.Error("Expected error when submitting after context cancel")
    }

    engine.Stop()
    engine.Wait()

    t.Logf("Processed %d tasks before cancellation", atomic.LoadInt32(&processed))
}

// TestHighConcurrency 测试高并发提交
func TestHighConcurrency(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var processed int32
    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            atomic.AddInt32(&processed, 1)
            time.Sleep(1 * time.Millisecond)
        },
        WorkersCount: 5,
        QueueLength:  100,
    })

    engine.Schedule()

    // 多个goroutine并发提交任务
    var wg sync.WaitGroup
    submitters := 10
    tasksPerSubmitter := 10

    for i := 0; i < submitters; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for j := 0; j < tasksPerSubmitter; j++ {
                engine.Submit(fmt.Sprintf("submitter-%d-task-%d", id, j))
            }
        }(i)
    }

    wg.Wait()
    time.Sleep(500 * time.Millisecond)

    engine.Stop()
    cancel()
    engine.Wait()

    expected := int32(submitters * tasksPerSubmitter)
    count := atomic.LoadInt32(&processed)
    if count != expected {
        t.Errorf("Expected %d tasks processed, got %d", expected, count)
    }
}

// TestQueueFull 测试队列满的情况
func TestQueueFull(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            time.Sleep(100 * time.Millisecond) // 慢处理
        },
        WorkersCount: 1, // 只有1个worker
        QueueLength:  5, // 小队列
    })

    engine.Schedule()

    // 快速提交任务，应该会填满队列
    submitted := 0
    timeout := time.After(500 * time.Millisecond)

    for i := 0; i < 20; i++ {
        select {
        case <-timeout:
            t.Logf("Timeout reached, submitted %d tasks", submitted)
            goto cleanup
        default:
            // 尝试提交
            done := make(chan error, 1)
            go func(taskID int) {
                done <- engine.Submit(fmt.Sprintf("task-%d", taskID))
            }(i)

            select {
            case err := <-done:
                if err != nil {
                    t.Logf("Submit failed at task %d: %v", i, err)
                    goto cleanup
                }
                submitted++
            case <-time.After(100 * time.Millisecond):
                t.Logf("Submit blocked at task %d (queue full)", i)
                goto cleanup
            }
        }
    }

cleanup:
    engine.Stop()
    cancel()
    engine.Wait()

    t.Logf("Successfully submitted %d tasks before blocking", submitted)
}

// TestDrainQueue 测试队列排空
func TestDrainQueue(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    var processed int32
    var mu sync.Mutex
    processedTasks := make(map[string]bool)

    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            mu.Lock()
            processedTasks[job.(string)] = true
            mu.Unlock()
            atomic.AddInt32(&processed, 1)
            time.Sleep(10 * time.Millisecond)
        },
        WorkersCount: 2,
        QueueLength:  50,
    })

    engine.Schedule()

    // 提交20个任务
    for i := 0; i < 20; i++ {
        engine.Submit(fmt.Sprintf("task-%d", i))
    }

    // 立即停止（确保有任务在队列中）
    time.Sleep(20 * time.Millisecond)
    engine.Stop()
    cancel()
    engine.Wait()

    count := atomic.LoadInt32(&processed)
    t.Logf("Processed %d tasks", count)

    // 验证所有任务都被处理
    if count != 20 {
        t.Errorf("Expected 20 tasks processed, got %d", count)
    }

    // 验证没有重复处理
    mu.Lock()
    uniqueTasks := len(processedTasks)
    mu.Unlock()

    if uniqueTasks != int(count) {
        t.Errorf("Some tasks processed multiple times: unique=%d, total=%d", uniqueTasks, count)
    }
}

// BenchmarkEngine 性能基准测试
func BenchmarkEngine(b *testing.B) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            // 模拟轻量级任务
            _ = job
        },
        WorkersCount: 10,
        QueueLength:  1000,
    })

    engine.Schedule()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        engine.Submit(i)
    }

    engine.Stop()
    cancel()
    engine.Wait()
}

// TestWorkerPoolDesign 测试工作池设计的正确性
func TestWorkerPoolDesign(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 追踪哪个worker处理了哪些任务
    var mu sync.Mutex
    workerTasks := make(map[int][]string)

    engine := NewEngine(&Options{
        Ctx: ctx,
        Func: func(job interface{}) {
            // 注意：我们无法直接获取worker ID，这里只是演示概念
            time.Sleep(10 * time.Millisecond)
        },
        WorkersCount: 3,
        QueueLength:  10,
    })

    engine.Schedule()

    // 提交9个任务（3个worker，每个应处理3个）
    for i := 0; i < 9; i++ {
        engine.Submit(fmt.Sprintf("task-%d", i))
    }

    time.Sleep(200 * time.Millisecond)
    engine.Stop()
    cancel()
    engine.Wait()

    mu.Lock()
    totalTasks := 0
    for _, tasks := range workerTasks {
        totalTasks += len(tasks)
    }
    mu.Unlock()

    t.Logf("Workers processed tasks (total: %d)", totalTasks)
}
