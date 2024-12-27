package main

import (
    "context"
    "log"
    "strconv"
    "testing"
    "time"
)

var Q *Queue

type StoreKafka struct{ opts *Options }

func NewStoreKafka(opts *Options) *StoreKafka { return &StoreKafka{opts: opts} }

func (sk *StoreKafka) Start() {
    Q = NewQueue(&Options{
        Ctx:          context.Background(),
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
        WorkersCount: 10,
        QueueLength:  1000,
    })
    sk.Start()

    for i := 0; i < 1000; i++ {
        Q.Chan() <- strconv.Itoa(i + 1)
    }

    time.Sleep(5 * time.Second)
    Q.Stop()
    cancel()
}
