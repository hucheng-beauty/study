package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 优雅关闭
func main() {
	var closing = make(chan struct{})
	var closed = make(chan struct{})

	go func() {
		// 模拟业务处理
		for {
			select {
			case <-closing:
				return
			default:
				// ....... 业务计算
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	// 处理 Ctrl + C 等中断信号
	quitChan := make(chan os.Signal)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
	<-quitChan

	close(closing)

	// 执行退出之前的清理动作
	go func(closed chan struct{}) {
		time.Sleep(time.Minute) // simulate cleanup
		close(closed)
	}(closed)

	select {
	case <-closed:
	case <-time.After(time.Second):
		fmt.Println("清理超时，不等了")
	}

	fmt.Println("graceful quit")
}
