package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arl/statsviz"
)

type HandleFunc func(job interface{})

type Options struct {
	JobChanLength int
	WorkerCount   int
	Fn            HandleFunc
}

type EfficientQueue struct {
	*Options
	jobs         []interface{}
	jobChans     []chan interface{}
	jobChan      chan interface{}
	jobChanQueue chan chan interface{}
}

func New(options *Options) *EfficientQueue {
	if options.JobChanLength <= 0 {
		options.JobChanLength = 1000
	}
	if options.WorkerCount <= 0 {
		options.WorkerCount = 100
	}

	return &EfficientQueue{
		Options:      options,
		jobChan:      make(chan interface{}, options.JobChanLength),
		jobChanQueue: make(chan chan interface{}, options.WorkerCount),
	}
}

func (eq *EfficientQueue) JobChan() chan interface{} {
	return eq.jobChan
}

func (eq *EfficientQueue) Run() {
	go eq.run()

	for i := 0; i < eq.WorkerCount; i++ {
		go eq.worker()
	}
}

func (eq *EfficientQueue) run() {
	var jobs []interface{}
	var jobChans []chan interface{}
	for {
		var activeJob interface{}
		var activeJobChan chan interface{}
		if len(jobs) > 0 && len(jobChans) > 0 {
			activeJob = jobs[0]
			activeJobChan = jobChans[0]
		}

		select {
		case job := <-eq.jobChan:
			jobs = append(jobs, job)
		case jobChan := <-eq.jobChanQueue:
			jobChans = append(jobChans, jobChan)
		case activeJobChan <- activeJob:
			jobs = jobs[1:]
			jobChans = jobChans[1:]
		}
	}
}

func (eq *EfficientQueue) worker() {
	for {
		// register the current worker into the worker pool
		eq.jobChanQueue <- eq.jobChan
		select {
		// deal with job
		case job := <-eq.jobChan:
			eq.Fn(job)
		}
	}
}

func main() {
	eq := New(&Options{
		JobChanLength: 1000,
		WorkerCount:   100,
		Fn: func(job interface{}) {
			log.Println(job)
		},
	})
	eq.Run()

	go func() {
		for {
			time.Sleep(100 * time.Microsecond)
			eq.JobChan() <- "hello world"
		}
	}()

	go func() {
		// runtime 图形化界面
		_ = statsviz.RegisterDefault()
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	// 处理 Ctrl + C 等中断信号
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	<-quitChan
	log.Println("main goroutine exit")
}
