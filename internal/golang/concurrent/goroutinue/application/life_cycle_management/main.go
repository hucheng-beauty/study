package main

import (
    "context"
    "fmt"
    "time"
)

type Tracker struct {
    ch   chan string
    stop chan struct{}
}

func NewTracker() *Tracker {
    return &Tracker{ch: make(chan string, 10)}
}

func (t *Tracker) Run() {
    for data := range t.ch {
        time.Sleep(time.Second * 1)
        fmt.Println(data)
    }
    t.stop <- struct{}{}
}

func (t *Tracker) Event(ctx context.Context, data string) error {
    select {
    case t.ch <- data:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (t *Tracker) Shutdown(ctx context.Context) {
    close(t.ch)
    select {
    case <-t.stop:
    case <-ctx.Done():
    }
}

// Never start a goroutine without knowing when it will stop
func main() {
    tr := NewTracker()

    // 上报 Leave concurrency to the caller
    go tr.Run()

    // 模拟埋点
    _ = tr.Event(context.Background(), "t1")
    _ = tr.Event(context.Background(), "t2")
    _ = tr.Event(context.Background(), "t3")

    time.Sleep(2 * time.Second)

    ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
    defer cancelFunc()

    tr.Shutdown(ctx)
}
