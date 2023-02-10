package patterns

import (
	"fmt"
	"time"
)

/*
	任务控制:
		非阻塞等待
		超时机制
		任务中断或退出
		优雅退出
*/

// ForSelect for select 无限循环模式
func ForSelect() {
	quit := make(chan int, 1)

	// send
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Printf("worker %d\n", i)
		}(i)

		//
		if i == 9 {
			quit <- i
		}
	}

	time.Sleep(time.Millisecond)

	// receive
	for {
		select {
		case <-quit:
			fmt.Println("over")
			return
		}
	}
}
