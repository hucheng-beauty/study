package queue_batch

import (
    "context"
    "log"
    "strconv"
    "testing"
    "time"
)

var Q *Queue

type QueueOptions struct {
    Ctx        context.Context
    BatchWait  int `json:"queue.batch.wait"`  // 批次等待时间(秒)
    BatchSize  int `json:"queue.batch.size"`  // 批次数量
    RoutineNum int `json:"queue.routine.num"` // 协程数量
    ChanSize   int `json:"queue.chan.size"`   // 管道缓存数量
}

type StoreKafka struct {
    queueInfo *QueueOptions
    // something ......
}

func NewStoreKafka(queueInfo *QueueOptions) *StoreKafka {
    return &StoreKafka{queueInfo: queueInfo}
}

func (sk *StoreKafka) Start() {
    Q = NewQueue(&Options{
        Ctx:        sk.queueInfo.Ctx,
        BatchWait:  time.Duration(sk.queueInfo.BatchWait) * time.Second,
        BatchSize:  sk.queueInfo.BatchSize,
        ChanSize:   sk.queueInfo.ChanSize,
        RoutineNum: sk.queueInfo.RoutineNum,
        Func:       sk.handle,
    })

    Q.Start()
}

func (sk *StoreKafka) handle(jobs []interface{}) {
    messages := make([]string, len(jobs))

    for i, job := range jobs {
        message, ok := job.(string)
        if !ok {
            continue
        }

        // pack message
        messages[i] = message
    }

    // send a message to kafka
    log.Println(messages)
}

func TestMainer(t *testing.T) {
    NewStoreKafka(&QueueOptions{
        Ctx:        context.Background(),
        BatchWait:  1,
        BatchSize:  200,
        RoutineNum: 5,
        ChanSize:   1000,
    }).Start()

    for i := 0; i < 1000; i++ {
        Q.Chan() <- strconv.Itoa(i + 1)
    }
    time.Sleep(3 * time.Second)
}
