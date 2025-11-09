package main

import (
    "fmt"
    "log"
    "time"
)

// 运行示例：go run stream_example.go run_example.go
func main() {
    // 运行流式接口示例
    RunStreamExample()

    fmt.Println("\n" + "="*60)
    fmt.Println("按回车键查看错误设计对比...")
    fmt.Scanln()

    // 运行错误设计示例
    WrongDesignExample()
}

// 如果只想运行其中一个示例，可以注释掉另一个
