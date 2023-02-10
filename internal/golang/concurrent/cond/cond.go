package main

import (
    "fmt"
    "sync"
    "time"
)

var (
    mutex sync.Mutex
    cond  = sync.NewCond(&mutex)
    queue []int
)

func produce(value int) {
    mutex.Lock()
    defer mutex.Unlock()

    queue = append(queue, value)
    fmt.Printf("Produced: %d\n", value)
    cond.Signal() // 唤醒一个等待的消费者
}

func consume() {
    mutex.Lock()
    defer mutex.Unlock()

    for len(queue) == 0 {
        cond.Wait() // 等待生产者生产
    }
    value := queue[0]
    queue = queue[1:]

    fmt.Printf("Consumed: %d\n", value)
}

func main() {
    for i := 0; i < 10; i++ {
        go produce(i)
        time.Sleep(time.Millisecond * 100) // 模拟生产速度
    }

    for i := 0; i < 10; i++ {
        go consume()
        time.Sleep(time.Millisecond * 150) // 模拟消费速度
    }

    // 等待所有协程完成
    time.Sleep(time.Second * 3)
}
