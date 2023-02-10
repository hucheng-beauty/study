package queue

import (
	"log"
	"testing"
	"time"
)

var Q *Queue

type QueueOptions struct {
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

func (s *StoreKafka) Start() {
	Q = NewQueue(&Options{
		BatchWait:  time.Duration(s.queueInfo.BatchWait) * time.Second,
		BatchSize:  s.queueInfo.BatchSize,
		ChanSize:   s.queueInfo.ChanSize,
		RoutineNum: s.queueInfo.RoutineNum,
		Func:       s.handle,
	})

	Q.Start()
}

type Message string

func (s *StoreKafka) handle(jobs []interface{}) {
	messages := make([]*Message, len(jobs))

	for i, job := range jobs {
		message, ok := job.(Message)
		if !ok {
			continue
		}

		// pack message
		messages[i] = &message
	}

	// send message to kafka
	log.Println(messages)
}

func TestMainer(t *testing.T) {
	NewStoreKafka(&QueueOptions{
		BatchWait:  1,
		BatchSize:  1000,
		RoutineNum: 100,
		ChanSize:   10000,
	}).Start()
}
