package main

import (
    "context"
    "log"
)

type HandleFunc func(job interface{})

type Options struct {
    Ctx          context.Context
    Func         HandleFunc
    WorkersCount int
    QueueLength  int
}

type Queue struct {
    workerPool chan chan interface{}
    jobChannel []chan interface{}
    jobQueue   chan interface{}
    quit       chan struct{}
    opts       *Options
}

func NewQueue(opts *Options) *Queue {
    return &Queue{
        workerPool: make(chan chan interface{}, opts.WorkersCount),
        jobChannel: make([]chan interface{}, opts.WorkersCount),
        jobQueue:   make(chan interface{}, opts.QueueLength),
        quit:       make(chan struct{}),
        opts:       opts,
    }
}

func (q *Queue) Start() {
    for i := 0; i < q.opts.WorkersCount; i++ {
        go func(i int) {
            q.jobChannel[i] = make(chan interface{})
            for {
                q.workerPool <- q.jobChannel[i]
                select {
                case job := <-q.jobChannel[i]:
                    q.opts.Func(job)
                case <-q.quit:
                    return
                }
            }
        }(i)
    }

    go func() {
        var jobQ []interface{}
        var workerQ []chan interface{}

        defer func() {
            if err := recover(); err != nil {
                log.Println("queue panic recover", err)
            }
        }()

        defer func() {
            for job := range q.jobQueue {
                jobQ = append(jobQ, job)
            }

            if len(jobQ) > 0 {
                for _, job := range jobQ {
                    q.opts.Func(job)
                }
            }
        }()

        for {
            var activeJob interface{}
            var activeWorker chan interface{}
            if len(jobQ) > 0 && len(workerQ) > 0 {
                activeJob = jobQ[0]
                activeWorker = workerQ[0]
            }

            select {
            case job := <-q.jobQueue:
                jobQ = append(jobQ, job)
            case worker := <-q.workerPool:
                workerQ = append(workerQ, worker)
            case activeWorker <- activeJob:
                jobQ = jobQ[1:]
                workerQ = workerQ[1:]
            case <-q.opts.Ctx.Done():
                return
            }
        }
    }()
}

func (q *Queue) Chan() chan interface{} { return q.jobQueue }

func (q *Queue) Stop() { close(q.quit) }
