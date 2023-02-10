package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
)

var (
    WorkerNumber  = 100
    JobChanLength = 1000
    JobChan       = make(chan interface{})
    JobChanQ      = make(chan chan interface{})
)

func main() {
    for i := 0; i < WorkerNumber; i++ {
        go func() {
            for {
                JobChanQ <- JobChan
                select {
                case job := <-JobChan:
                    log.Println(job)
                }
            }
        }()
    }

    // send
    go func() {
        for i := 0; i < 10000; i++ {
            // time.Sleep(time.Millisecond * 100)
            JobChan <- "hello world"
        }
    }()

    go func() {
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
            case job := <-JobChan:
                jobs = append(jobs, job)
            case jobChan := <-JobChanQ:
                jobChans = append(jobChans, jobChan)
            case activeJobChan <- activeJob:
                jobs = jobs[1:]
                jobChans = jobChans[1:]
            }
        }
    }()

    // 处理 Ctrl + C 等中断信号
    quitChan := make(chan os.Signal)
    signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
    <-quitChan
    log.Println("main goroutine exit")
}
