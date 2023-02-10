package main

import (
    "fmt"
    "strings"
)

func collect(ip string, resultChan chan string, limitChan chan struct{}) {
    defer func() {
        resultChan <- strings.ToUpper(ip)
        <-limitChan
        close(resultChan)
    }()
}

func main() {
    ips := []string{"c", "cpp", "go", "python", "java", "rust", "ruby"}
    chs := make([]chan string, len(ips)) // 采集信息通过 channel 收集
    limitChan := make(chan struct{}, 3)  // 控制 goroutine 的数量,达到限流的效果

    for i, ip := range ips {
        chs[i] = make(chan string, 1)
        limitChan <- struct{}{}

        go collect(ip, chs[i], limitChan)
    }

    for _, ch := range chs {
        for v := range ch {
            fmt.Println(v) // 处理采集信息
        }
    }
}
