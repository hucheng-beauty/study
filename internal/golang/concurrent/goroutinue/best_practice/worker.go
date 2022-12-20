package main

import (
	"log"
	"time"
)

type Worker struct {
	WorkerPool chan chan string
	JobChannel chan string
	quit       chan bool
}

func NewWorker(workerPool chan chan string) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan string),
		quit:       make(chan bool)}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel // register the current worker into the worker queue.
			select {
			case job := <-w.JobChannel:
				// deal with job
				log.Println(time.Now().UnixMilli(), job)
			case <-w.quit:
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
