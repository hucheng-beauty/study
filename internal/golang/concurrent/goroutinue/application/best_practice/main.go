package main

import (
	"os"
	"strconv"
	"time"
)

/*
	各个组件:
		JobChan
			对请求进行缓存

		Dispatcher
			jobChannel(chan string) <- WorkerPool(chan chan string)
			jobChannel(chan string) <- job(string)

		Worker
			WorkerPool(chan chan string) <- jobChannel(chan string)
			job(string) <- jobChannel(chan string)

	整体数据流向					Dispatcher			WorkerPool                  Worker

		   request 									JobChannel					worker
		   request		==》	JobChan	==》	JobChannel		==》		worker
		   request 									JobChannel					worker
*/

var (
	MaxWorker = os.Getenv("MAX_WORKERS")
	maxWorker int

	MaxQueue = os.Getenv("MAX_QUEUE")
	maxQueue int
)

var JobQueue = make(chan string, maxQueue)

func main() {
	mw, err := strconv.Atoi(MaxWorker)
	if err != nil {
		maxWorker = 10
	}
	maxWorker = mw

	mq, err := strconv.Atoi(MaxQueue)
	if err != nil {
		maxQueue = 10
	}
	maxQueue = mq

	NewDispatcher(maxWorker).Run()

	go func() {
		for {
			JobQueue <- "hello"
		}
	}()
	time.Sleep(time.Millisecond)
}
