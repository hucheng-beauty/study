package main

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

/*
	SingleFlight 是 Go 开发组提供的一个扩展并发原语;
	它的作用是在处理多个 goroutine 同时调用同一个函数的时候,
	只让一个 goroutine 去调用这个函数,等到这个 goroutine 返回结果的时候,
	再把结果返回给这几个同时调用的 goroutine,这样可以减少并发调用的数量
*/

type H2O struct {
	semaH *semaphore.Weighted
	semaO *semaphore.Weighted
	wg    sync.WaitGroup // 将循环栅栏替换成WaitGroup
}

func New() *H2O {
	var wg sync.WaitGroup
	wg.Add(3)

	return &H2O{
		semaH: semaphore.NewWeighted(2),
		semaO: semaphore.NewWeighted(1),
		wg:    wg,
	}
}

func (h2o *H2O) hydrogen(releaseHydrogen func()) {
	h2o.semaH.Acquire(context.Background(), 1)
	releaseHydrogen()

	// 标记自己已达到，等待其它goroutine到达
	h2o.wg.Done()
	h2o.wg.Wait()

	h2o.semaH.Release(1)
}

func (h2o *H2O) oxygen(releaseOxygen func()) {
	h2o.semaO.Acquire(context.Background(), 1)
	releaseOxygen()

	// 标记自己已达到，等待其它goroutine到达
	h2o.wg.Done()
	h2o.wg.Wait()
	// 都到达后重置 wg
	h2o.wg.Add(3)

	h2o.semaO.Release(1)
}

func main() {

}
