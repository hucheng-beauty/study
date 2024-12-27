package main

import (
    "fmt"

    "github.com/HdrHistogram/hdrhistogram-go"
)

func main() {
    h := hdrhistogram.New(1, 100000, 3) // 范围为 1 到 100000，精度为 3 位小数。

    // 模拟记录响应时间（以毫秒为单位）。
    h.RecordValue(100)
    h.RecordValue(200)
    h.RecordValue(300)

    // 计算 TP99
    tp99 := h.ValueAtQuantile(99)
    fmt.Printf("TP99: %dms\n", tp99)
}
