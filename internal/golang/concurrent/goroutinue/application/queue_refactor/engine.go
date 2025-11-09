package queue_refactor

import (
    "context"
    "log"
    "sync"
)

type HandleFunc func(job interface{})

type Options struct {
    Ctx          context.Context
    Func         HandleFunc
    WorkersCount int
    QueueLength  int
}

type Engine struct {
    wg         sync.WaitGroup
    once       sync.Once
    opts       *Options
    workerPool chan chan interface{}
    jobChannel []chan interface{}
    jobQueue   chan interface{}
    done       chan struct{}
}

func NewEngine(opts *Options) *Engine {
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
        opts.QueueLength = 10000
    }

    q := &Engine{
        workerPool: make(chan chan interface{}, opts.WorkersCount),
        jobChannel: make([]chan interface{}, opts.WorkersCount),
        jobQueue:   make(chan interface{}, opts.QueueLength),
        done:       make(chan struct{}),
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

func (e *Engine) handle(job interface{}) {
    if job == nil {
        return
    }

    log.Println("Handling jobs:", job)
}

// Submit submits a job to the queue
func (e *Engine) Submit(job interface{}) error {
    select {
    case e.jobQueue <- job:
        return nil
    case <-e.opts.Ctx.Done():
        return e.opts.Ctx.Err()
    }
}

func (e *Engine) Schedule() error {
    // start workers
    for i := 0; i < e.opts.WorkersCount; i++ {
        e.wg.Add(1)
        go func(i int) {
            defer func() {
                if err := recover(); err != nil {
                    log.Printf("worker %d panic recover: %v", i, err)
                }
            }()

            e.worker(i)
        }(i)
    }

    // engine begin to schedule
    e.wg.Add(1)
    go func() {
        defer e.wg.Done()
        defer func() {
            if err := recover(); err != nil {
                log.Println("engine dispatch panic recover:", err)
            }
        }()

        e.schedule()
    }()

    return nil
}

func (e *Engine) worker(i int) {
    defer e.wg.Done()

    for {
        select {
        case e.workerPool <- e.jobChannel[i]:
            // successfully register to the working pool and wait for the task
            select {
            case job := <-e.jobChannel[i]:
                e.opts.Func(job)
            case <-e.opts.Ctx.Done():
                return
            }

        case <-e.opts.Ctx.Done():
            return
        }
    }
}

func (e *Engine) schedule() {
    for {
        select {
        case job := <-e.jobQueue:
            select {
            case jobChannel := <-e.workerPool:
                select {
                // submit job to jobChannel
                case jobChannel <- job:

                // context cancels, and exits after processing the current task
                case <-e.opts.Ctx.Done():
                    e.opts.Func(job)
                    return
                }
            }
        case <-e.done:
            e.drainQueue()
            return
        }
    }
}

func (e *Engine) drainQueue() {
    log.Println("draining remaining jobs in queue...")

    drained := 0

    for {
        select {
        case job, ok := <-e.jobQueue:
            if !ok {
                // jobQueue has closed
                log.Printf("drained %d jobs, queue closed\n", drained)
                return
            }

            // handle job directly, no longer assign to workers,
            // because of  workers may be exiting
            e.opts.Func(job)
            drained++
        default:
            // queue is empty, exit
            log.Printf("drained %d jobs, queue is empty\n", drained)
            return
        }
    }
}

// Stop tells the scheduler to stop new job(processing the remaining jobs),
// but does not wait for the workers to quit.
// Note: cancel of context should be called immediately after calling
// Stop to inform workers to quit, and then Wait should be called.
func (e *Engine) Stop() {
    e.once.Do(func() {
        log.Println("engine stop...")
        // close done channel to stop scheduling new jobs
        close(e.done)
    })
}

// Wait is wait for all workers to exit
func (e *Engine) Wait() { e.wg.Wait() }
