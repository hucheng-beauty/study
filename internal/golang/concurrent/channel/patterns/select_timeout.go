package patterns

import (
	"fmt"
	"time"
)

// SelectTimeout select timeout 模式
func SelectTimeout() {
	ch := make(chan string)
	timeout := time.After(3 * time.Second)

	// send
	go func() {
		// 模拟网络访问
		time.Sleep(5 * time.Second)
		ch <- "网络访问阻塞"
	}()

	// receive
	for {

		select {
		case v := <-ch:
			fmt.Println(v)
		case <-timeout:
			fmt.Println("超时退出")
			return
		default:
			fmt.Println("等待...")
			time.Sleep(1 * time.Second)
		}
	}
}
