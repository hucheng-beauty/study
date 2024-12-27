package main

import (
    "context"
    "fmt"
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
    // compliance checks
    if opts == nil {
        opts = &Options{}
    }
    if opts.Ctx == nil {
        opts.Ctx = context.Background()
    }
    if opts.WorkersCount <= 0 {
        opts.WorkersCount = 10
    }
    if opts.QueueLength <= 0 {
        opts.QueueLength = 1000
    }

    q := &Queue{
        workerPool: make(chan chan interface{}, opts.WorkersCount),
        jobChannel: make([]chan interface{}, opts.WorkersCount),
        jobQueue:   make(chan interface{}, opts.QueueLength),
        quit:       make(chan struct{}),
        opts:       opts,
    }

    if opts.Func == nil {
        opts.Func = q.handle
    }

    // initialize jobChannel's item
    for i := 0; i < opts.WorkersCount; i++ {
        q.jobChannel[i] = make(chan interface{})
    }

    return q
}

func (q *Queue) Chan() chan interface{} { return q.jobQueue }

func (q *Queue) handle(job interface{}) {
    if job == nil {
        return
    }
    fmt.Println("handle job:", job)
}

func (q *Queue) worker(i int) {
    defer func() {
        for job := range q.jobQueue {
            q.opts.Func(job)
        }
    }()

    for {
        // register q.jobChannel[i] to q.workerPool
        q.workerPool <- q.jobChannel[i]
        select {
        case job := <-q.jobChannel[i]:
            q.opts.Func(job)
        case <-q.quit:
            return
        }
    }
}

func (q *Queue) Start() {
    for i := 0; i < q.opts.WorkersCount; i++ {
        go func() {
            defer func() {
                if err := recover(); err != nil {
                    log.Println("worker panic recover", err)
                }
            }()
            q.worker(i)
        }()
    }

    go func() {
        defer func() {
            if err := recover(); err != nil {
                log.Println("queue panic recover", err)
            }
        }()
        q.start()
    }()
}

func (q *Queue) start() {
    var jobQ []interface{}
    var workerQ []chan interface{}

    defer func() {
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
}
