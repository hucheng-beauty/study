package main

import (
	"fmt"
	"time"
)

/*
	协程 Coroutine
		轻量级 "线程"
		"非抢占式多任务处理",由协程主动交出控制权
		编译器/解释器/虚拟机层面的多任务
		多个协程可能在一个或多个线程上运行

	goroutine 的定义
		任何函数只需加上 go 就能送给调度器运行
		不需要在定义时区分是否为异步函数
		调度器在合适的点进行切换
		使用 -race 来检查数据访问冲突
			eg. go run -race file_name.go

	goroutine 可能切换点(只做参考,不能保证切换)
		I/O, select
		channel
		等待🔒
		函数调用时
		runtime.Gosched()
*/

func test() {
	var arr [10]int
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				arr[i]++ // write
				// runtime.Gosched() // 手动主动交出控制权
			}
		}(i)
	}

	time.Sleep(time.Millisecond)
	fmt.Println(arr) // read
}
