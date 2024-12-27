package queue

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
    Quit       chan struct{}
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
        Quit:       make(chan struct{}),
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
    for {
        // register q.jobChannel[i] to q.workerPool
        q.workerPool <- q.jobChannel[i]
        select {
        case job := <-q.jobChannel[i]:
            fmt.Println("[worker1]", i)
            q.opts.Func(job)
        case <-q.opts.Ctx.Done():
            fmt.Println("[worker] receive Done:", i)
            return
        }
    }
}

func (q *Queue) Start() {
    for i := 0; i < q.opts.WorkersCount; i++ {
        go func(j int) {
            defer func() {
                if err := recover(); err != nil {
                    log.Println("worker panic recover", err)
                }
            }()
            q.worker(j)
        }(i)
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
    defer func() {
        for job := range q.jobQueue {
            q.opts.Func(job)
        }
        fmt.Println("[start] queue quit:")
    }()

    for {
        select {
        case job := <-q.jobQueue:
            jobChannel := <-q.workerPool
            jobChannel <- job
        case <-q.Quit:
            fmt.Println("[start] receive quit:")
            return
        }
    }
}
