package engine

import "log"

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

// ConcurrentEngine :concurrent version engine.
type ConcurrentEngine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ItemChan         chan Item
	RequestProcessor Processor
}

type Processor func(Request) (ParseResult, error)

func (e *ConcurrentEngine) Run(seeds ...Request) {
	// start Scheduler and create requestChan and queueChan
	e.Scheduler.Run()

	// create worker and wait for request into requestChan
	out := make(chan ParseResult)
	for i := 0; i < e.WorkerCount; i++ {
		in := e.Scheduler.WorkerChan()
		e.createWorker(in, out, e.Scheduler)
	}

	// submit request to requestChan
	for _, r := range seeds {
		e.Scheduler.Submit(r)
	}

	for {
		result := <-out
		// send item to ItemChan and storage item
		for _, item := range result.Items {
			go func(i Item) { e.ItemChan <- i }(item)
		}

		// submit request to requestChan
		for _, request := range result.Requests {
			// URL De-duplication
			if !isDuplicate(request.Url) {
				log.Printf("repeated url: %v", request.Url)
				continue
			}
			e.Scheduler.Submit(request)
		}
	}
}

func (e *ConcurrentEngine) createWorker(in chan Request, out chan ParseResult, ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in
			result, err := e.RequestProcessor(request)
			if err != nil {
				continue
			}
			/*
				// call rpc
				result, err := Worker(request) // worker start work
				if err != nil {
					continue
				}
			*/
			out <- result
		}
	}()
}

// TODO: restart is not storage data
var visitedUrls = make(map[string]bool, 0)

func isDuplicate(url string) bool {
	if visitedUrls[url] {
		return true
	}

	visitedUrls[url] = true
	return false
}
