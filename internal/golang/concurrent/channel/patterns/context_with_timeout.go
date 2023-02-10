package patterns

import (
	"context"
	"fmt"
	"time"
)

// ContextWithTimeout Context 的 WithTimeout 函数超时取消
func ContextWithTimeout() {
	// 创建一个子节点的context,3秒后自动超时
	ctx, exit := context.WithTimeout(context.Background(), 3*time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("1 下班咯~~~")
				return
			default: // select 中有 default 则为非阻塞等待
				fmt.Println("worker bee 1", "认真摸鱼中，请勿打扰...")
			}
			time.Sleep(1 * time.Second)
		}

	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("2 下班咯~~~")
				return
			default:
				fmt.Println("worker bee 2", "认真摸鱼中，请勿打扰...")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	time.Sleep(5 * time.Second) // 工作5秒后休息,5秒后发出停止指令
	exit()
	fmt.Println("结束了")
}
