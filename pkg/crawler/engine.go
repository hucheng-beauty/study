package crawler

import (
    "log"

    "study/internal/crawler/fetcher"
)

type Scheduler interface {
    Run()
    Submit(Request)
    JobChan() chan Request
    Dispatcher(chan Request)
}

type Processor func(Request) (Result, error)

func FetchProcessor(r Request) (Result, error) {
    body, err := fetcher.Fetch(r.Url)
    if err != nil {
        return Result{}, err
    }
    return r.Parser.Parse(body, r.Url), nil
}

type Engine struct {
    Scheduler        Scheduler
    RequestProcessor Processor
    WorkerCount      int
    ItemChan         chan Item
}

func (e *Engine) Run(seeds ...Request) {
    go e.Scheduler.Run()

    out := make(chan Result)
    for i := 0; i < e.WorkerCount; i++ {
        // 1.dispatcher job chan into job chan queue, jobChanQueue <- jobChan
        // 2.deal with job
        go func() {
            for {
                jobChan := e.Scheduler.JobChan()
                e.Scheduler.Dispatcher(jobChan)
                select {
                case job := <-jobChan:
                    result, err := e.RequestProcessor(job)
                    if err != nil {
                        log.Println(err)
                        continue
                    }
                    out <- result
                }
            }
        }()
    }

    // send request
    for _, r := range seeds {
        e.Scheduler.Submit(r)
    }

    // deal with result and continue submitting request
    for {
        result := <-out

        // deal with item
        for _, item := range result.Items {
            go func(i Item) { e.ItemChan <- i }(item)
        }

        // deal with result
        for _, request := range result.Requests {
            if !isDuplicate(request.Url) {
                log.Printf("repeated url: %v", request.Url)
                continue
            }
            e.Scheduler.Submit(request)
        }
    }
}