package queue

import (
	"context"
	"log"
	"time"
)

type HandleFunc func(jobs []interface{})

type Options struct {
	Ctx        context.Context
	BatchWait  time.Duration
	BatchSize  int
	ChanSize   int
	RoutineNum int
	Func       HandleFunc
}

type Queue struct {
	c chan interface{}
	*Options
}

func NewQueue(options *Options) *Queue {
	if options.BatchWait == 0 {
		options.BatchWait = time.Second
	}
	if options.BatchSize == 0 {
		options.BatchSize = 100
	}
	if options.ChanSize == 0 {
		options.ChanSize = 100
	}
	if options.RoutineNum == 0 {
		options.RoutineNum = 100
	}

	return &Queue{
		c:       make(chan interface{}, options.ChanSize),
		Options: options,
	}
}

func (q *Queue) Chan() chan interface{} {
	return q.c
}

func (q *Queue) Start() {
	for i := 0; i < q.RoutineNum; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("queue panic recover", err)
				}
			}()
			q.start()
		}()
	}
}

func (q *Queue) start() {
	jobs := make([]interface{}, 0)

	defer func() {
		for job := range q.c {
			jobs = append(jobs, job)
		}
		q.Func(jobs)
	}()

	ticker := time.NewTicker(q.BatchWait)
	for {
		select {
		case job, ok := <-q.c:
			if !ok {
				continue
			}

			jobs = append(jobs, job)
			if len(jobs) < q.BatchSize {
				continue
			}

			goto LOOP
		case <-ticker.C:
			if len(jobs) == 0 {
				continue
			}

			goto LOOP
		case <-q.Ctx.Done():
			ticker.Stop()
			return
		}
	}
LOOP:
	q.Func(jobs)
	jobs = make([]interface{}, 0)
}
