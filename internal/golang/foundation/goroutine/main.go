package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

func worker(ctx context.Context) {
	defer wg.Done()
LOOP:
	for {
		fmt.Println("working...")
		time.Sleep(time.Second)
		select {
		case <-ctx.Done():
			break LOOP
		default:
		}
	}
}

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	wg.Add(1)
	go worker(ctx)

	time.Sleep(3 * time.Second)
	cancelFunc()
	wg.Wait()
}
