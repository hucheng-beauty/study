package main

import "fmt"

type ReadyNotifier interface {
	WorkerReady(chan string)
}

type Scheduler interface {
	ReadyNotifier
	Submit(string)
	WorkerChan() chan string
	Run()
}

type Engine struct {
	Scheduler        Scheduler
	WorkerCount      int
	ResultChan       chan string
	RequestProcessor Processor
}

type Processor func(str string) (string, error)

func (e *Engine) Run(seeds ...string) {
	e.Scheduler.Run()

	out := make(chan string)
	for i := 0; i < e.WorkerCount; i++ {
		go func(in chan string, out chan string, ready ReadyNotifier) {
			// chan chan string <- chan string
			ready.WorkerReady(in)

			// <- chan string
			request := <-in
			r, err := e.RequestProcessor(request)
			if err != nil {
				return
			}
			// chan string <- string
			out <- r
		}(e.Scheduler.WorkerChan(), out, e.Scheduler)
	}

	// submit request to requestChan
	for _, r := range seeds {
		// chan string <- string
		e.Scheduler.Submit(r)
	}

	for {
		r := <-out
		fmt.Println(r)
	}
}

func main() {
	// config Engine
	e := Engine{
		Scheduler:        &QueuedScheduler{},
		WorkerCount:      5,
		ResultChan:       nil,
		RequestProcessor: Worker,
	}

	// start run
	e.Run("c", "python", "java", "go", "cpp", "ruby")
	select {}
}
