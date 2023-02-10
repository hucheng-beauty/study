package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	limitChan := make(chan struct{}, 2)
	ips := []string{"c", "cpp", "go", "python", "java", "rust", "ruby"}

	for _, ip := range ips {
		wg.Add(1)
		limitChan <- struct{}{}

		go func(ip string) {
			doneChan := make(chan bool)

			ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancelFunc()

			go worker(ctx, &wg, ip, doneChan)

			select {
			case <-ctx.Done():
				fmt.Println(ip, " 超时了!")
				<-limitChan
				wg.Done()

			case done := <-doneChan:
				if done {
					fmt.Println(ip, " 正常完成!")
					<-limitChan
					wg.Done()
				}
			}
		}(ip)
	}
	wg.Wait()
}

func worker(ctx context.Context, wg *sync.WaitGroup, name string, taskDoneChan chan bool) {
	defer func() {
		taskDoneChan <- true
		fmt.Println(name, " 处理完成并释放资源!")
	}()

	randInt := rand.Intn(5)
	fmt.Println(name, "==>", randInt)
	time.Sleep(time.Second * time.Duration(randInt))
}
