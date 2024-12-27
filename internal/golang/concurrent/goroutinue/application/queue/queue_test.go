package queue

import (
    "context"
    "fmt"
    "log"
    "runtime"
    "strconv"
    "testing"
    "time"
)

var Q *Queue

type StoreKafka struct{ opts *Options }

func NewStoreKafka(opts *Options) *StoreKafka { return &StoreKafka{opts: opts} }

func (sk *StoreKafka) Start() {
    Q = NewQueue(&Options{
        Ctx:          sk.opts.Ctx,
        Func:         sk.handle,
        WorkersCount: sk.opts.WorkersCount,
        QueueLength:  sk.opts.QueueLength,
    })

    Q.Start()
}

func (sk *StoreKafka) handle(job interface{}) {
    if job == nil {
        return
    }
    log.Println("[handle] job: ", job)
}

func TestQueue(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    sk := NewStoreKafka(&Options{
        Ctx:          ctx,
        WorkersCount: 100,
        QueueLength:  1000,
    })

    sk.Start()

    for i := 0; i < 1000; i++ {
        Q.Chan() <- strconv.Itoa(i + 1)
    }

    time.Sleep(time.Second * 3)
    fmt.Println("go number 1:", runtime.NumGoroutine())
    Q.Quit <- struct{}{}
    cancel()
    close(Q.Chan())
    time.Sleep(time.Second * 1)
    fmt.Println("go number 2:", runtime.NumGoroutine())
}
