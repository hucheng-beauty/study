package toolkit

import (
    "context"
    "log"
    "sync"
    "time"
)

type HandleFunc func(jobs []interface{})

type Options struct {
    Ctx          context.Context // context for cancellation
    Func         HandleFunc      // user-defined batch processing function
    BatchWait    time.Duration   // maximum wait time for each batch
    BatchSize    int             // maximum number of jobs per batch
    ChanSize     int             // buffer size of the internal channel
    RoutineNum   int             // number of worker goroutines
    DrainTimeout time.Duration   // timeout for draining remaining jobs when context is canceled
}

type Queue struct {
    c         chan interface{}
    wg        sync.WaitGroup
    closeOnce sync.Once
    *Options
}

func NewQueue(options *Options) *Queue {
    if options.BatchWait <= 0 {
        options.BatchWait = 60 * time.Second
    }
    if options.BatchSize <= 0 {
        options.BatchSize = 100
    }
    if options.ChanSize <= 0 {
        options.ChanSize = 1000000
    }
    if options.RoutineNum <= 0 {
        options.RoutineNum = 10
    }
    if options.DrainTimeout <= 0 {
        options.DrainTimeout = 5 * time.Second
    }

    return &Queue{c: make(chan interface{}, options.ChanSize), Options: options}
}

func (q *Queue) Chan() chan interface{} { return q.c }

func (q *Queue) Start() {
    for i := 0; i < q.RoutineNum; i++ {
        q.wg.Add(1)
        go func() {
            defer func() {
                q.wg.Done()

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
    ticker := time.NewTicker(q.BatchWait)

    defer func() {
        ticker.Stop()
        // process current batch if any
        if len(jobs) > 0 {
            q.safeFunc(jobs)
            jobs = jobs[:0]
        }
    }()

    for {
        select {
        case job, ok := <-q.c:
            if !ok {
                return // when channel closed, exit goroutine
            }

            jobs = append(jobs, job)
            if len(jobs) < q.BatchSize {
                continue
            }

            goto LOOP
        case <-ticker.C:
            if len(jobs) <= 0 {
                continue
            }

            goto LOOP
        case <-q.Ctx.Done():
            // when context canceled, try to drain remaining jobs with timeout
            q.drainWithTimeout(jobs)

            // clear jobs to avoid duplicate processing in defer
            jobs = jobs[:0]
            return
        }
    LOOP:
        q.safeFunc(jobs)
        jobs = make([]interface{}, 0)
    }
}

// drainWithTimeout processes remaining jobs in the channel with a timeout
func (q *Queue) drainWithTimeout(currentBatch []interface{}) {
    jobs := currentBatch
    timeout := time.NewTimer(q.DrainTimeout)
    defer timeout.Stop()

    // process current batch first
    if len(jobs) > 0 {
        q.safeFunc(jobs)
        jobs = jobs[:0]
    }

    // try to drain channel with timeout
    for {
        select {
        case job, ok := <-q.c:
            if !ok {
                // channel closed
                if len(jobs) > 0 {
                    q.safeFunc(jobs)
                }
                return
            }

            jobs = append(jobs, job)
            if len(jobs) >= q.BatchSize {
                q.safeFunc(jobs)
                jobs = jobs[:0]
            }

        case <-timeout.C:
            // timeout reached, process remaining and exit
            if len(jobs) > 0 {
                q.safeFunc(jobs)
            }
            return

        default:
            // channel empty, process remaining and exit
            if len(jobs) > 0 {
                q.safeFunc(jobs)
            }
            return
        }
    }
}

// safeFunc wraps the user's handle function with panic recovery
func (q *Queue) safeFunc(jobs []interface{}) {
    defer func() {
        if err := recover(); err != nil {
            log.Printf("queue handle panic: %v, dropped %d jobs\n", err, len(jobs))
        }
    }()

    q.Func(jobs)
}

// Close gracefully shuts down the queue
func (q *Queue) Close() {
    q.closeOnce.Do(func() { close(q.c) })
    q.wg.Wait()
}
