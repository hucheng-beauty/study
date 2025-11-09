package toolkit

import (
    "context"
    "log"
    "strconv"
    "testing"
    "time"
)

var StoreKafkaQ *Queue

type StoreKafka struct{ opts *Options }

func NewStoreKafka(opts *Options) *StoreKafka { return &StoreKafka{opts: opts} }

func (sk *StoreKafka) Start() {
    StoreKafkaQ = NewQueue(&Options{
        Ctx:        context.Background(),
        BatchWait:  time.Duration(sk.opts.BatchWait) * time.Second,
        BatchSize:  sk.opts.BatchSize,
        ChanSize:   sk.opts.ChanSize,
        RoutineNum: sk.opts.RoutineNum,
        Func:       sk.handle,
    })

    StoreKafkaQ.Start()
}

func (sk *StoreKafka) handle(jobs []interface{}) {
    log.Printf("[handle] jobs:%v\n", jobs)
    if jobs == nil || 0 >= len(jobs) {
        return
    }

    messages := make([]string, len(jobs))
    for i, job := range jobs {
        message, ok := job.(string)
        if !ok {
            continue
        }
        messages[i] = message
    }
    log.Println("[handle] messages:", messages)
}

func TestQueue(t *testing.T) {
    sk := NewStoreKafka(&Options{
        Ctx:        context.Background(),
        BatchWait:  2,
        BatchSize:  200,
        RoutineNum: 1,
        ChanSize:   200,
    })
    sk.Start()

    for i := 0; i < 1000; i++ {
        StoreKafkaQ.Chan() <- strconv.Itoa(i)
    }
}
