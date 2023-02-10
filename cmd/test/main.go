// package main
//
// import (
//     "fmt"
//     "log"
//     "os"
//     "os/signal"
//     "sync"
//     "syscall"
//     "time"
// )
//
// var (
//     WorkerNumber  = 100
//     JobChanLength = 1000
//     JobChan       = make(chan interface{})
//     JobChanQ      = make(chan chan interface{})
// )
//
// func main() {
//     for i := 0; i < WorkerNumber; i++ {
//         go func() {
//             for {
//                 JobChanQ <- JobChan
//                 select {
//                 case job := <-JobChan:
//                     log.Println(job)
//                 }
//             }
//         }()
//     }
//
//     // send
//     go func() {
//         for i := 0; i < 10000; i++ {
//             // time.Sleep(time.Millisecond * 100)
//             JobChan <- "hello world"
//         }
//     }()
//
//     go func() {
//         var jobs []interface{}
//         var jobChans []chan interface{}
//         for {
//             var activeJob interface{}
//             var activeJobChan chan interface{}
//             if len(jobs) > 0 && len(jobChans) > 0 {
//                 activeJob = jobs[0]
//                 activeJobChan = jobChans[0]
//             }
//
//             select {
//             case job := <-JobChan:
//                 jobs = append(jobs, job)
//             case jobChan := <-JobChanQ:
//                 jobChans = append(jobChans, jobChan)
//             case activeJobChan <- activeJob:
//                 jobs = jobs[1:]
//                 jobChans = jobChans[1:]
//             }
//         }
//     }()
//
//     // 处理 Ctrl + C 等中断信号
//     quitChan := make(chan os.Signal)
//     signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
//     <-quitChan
//     log.Println("main goroutine exit")
// }

package main

import (
    "fmt"
    "sync"
)

// POOL 定义打印范围
var POOL = 100

// printOdd 函数用于打印奇数
func printOdd(p chan int, group *sync.WaitGroup) {
    defer group.Done()
    // 输出奇数
    for i := 1; i <= POOL; i++ {
        // 将 i 发送到通道 p
        p <- i //
        if i%2 == 1 {
            fmt.Println("奇数:", i) //
        }
    }
}

// printEven 函数用于打印偶数
func printEven(p chan int, group *sync.WaitGroup) {
    defer group.Done()
    for i := 1; i <= POOL; i++ {
        // 从通道 p 接收数据，并赋值给 num
        num := <-p
        // time.Sleep(100 * time.Millisecond)
        if num%2 == 0 {
            fmt.Println("偶数:", num)
        }
    }
}

func main() {
    // 创建一个无缓冲的整型通道
    msg := make(chan int, 1)
    var s sync.WaitGroup
    s.Add(2)
    // 启动两个 goroutine
    go printOdd(msg, &s)
    go printEven(msg, &s)
    s.Wait()
}
