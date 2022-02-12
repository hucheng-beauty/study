package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	println("main goroutine:", ctx.Value("hello"))
	withValue := context.WithValue(ctx, "hello", "world")

	go worker(withValue)

	time.Sleep(2 * time.Second)
	println("main goroutine will do cancelFunc")
	cancelFunc()
}

func worker(ctx context.Context) {
	fmt.Println("worker goroutine:", ctx.Value("hello"))
	cancel, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	println("worker starting")
	time.Sleep(10 * time.Second)
	println("worker ending")

	select {
	case <-cancel.Done():
		return
	}
}
