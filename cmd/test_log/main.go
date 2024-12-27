package main

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
)

const (
    Green = "\033[32m"
    Red   = "\033[31m"
    Reset = "\033[0m"
)

func main() {
    // 获取当前执行目录
    // execDir, err := os.Getwd()
    // if err != nil {
    //     log.Fatalf("获取当前目录失败: %v", err)
    // }

    // 拼接日志文件路径
    logFilePath := filepath.Join(".", "app.log")

    // 打开日志文件，如果不存在则创建，追加写入
    logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("打开日志文件失败: %v", err)
    }
    defer logFile.Close()

    // 设置日志输出到文件
    log.SetOutput(logFile)
    log.Default()

    // 写入日志
    log.Printf("%s这是一条日志信息%s", Green, Reset)
    // log.Printf("当前执行目录: %s", execDir)

    fmt.Println("日志已写入到文件:", logFilePath)
}
