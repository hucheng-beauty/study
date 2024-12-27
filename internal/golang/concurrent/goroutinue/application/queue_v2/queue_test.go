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
    NewStoreKafka(&Options{
        Ctx:          context.Background(),
        WorkersCount: 100,
        QueueLength:  1000,
    }).Start()

    for i := 0; i < 1000; i++ {
        Q.Chan() <- strconv.Itoa(i + 1)
    }

    time.Sleep(3 * time.Second)
}
