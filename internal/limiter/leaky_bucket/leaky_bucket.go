package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Result 封装业务逻辑的执行结果
type Result struct {
	Msg string
}

// Handler 执行的业务逻辑函数
type Handler func() Result

// Task 每个请求来了,把需要执行的业务逻辑封装成 Task，放入木桶,等待 worker 取出执行
type Task struct {
	handler    Handler     // worker 从木桶中取出请求对象后要执行的业务逻辑函数
	resultChan chan Result // 等待 worker 执行并返回结果的 channel
	taskID     int
}

func NewTask(id int, handler Handler) Task {
	return Task{
		handler:    handler,
		resultChan: make(chan Result),
		taskID:     id,
	}
}

type LeakyBucket struct {
	BucketSize int       // 木桶的大小
	WorkerNum  int       // 同时从木桶中获取任务执行的 worker 数量
	bucket     chan Task // 存放任务的木桶
}

func NewLeakyBucket(bucketSize int, workNum int) *LeakyBucket {
	return &LeakyBucket{
		BucketSize: bucketSize,
		WorkerNum:  workNum,
		bucket:     make(chan Task, bucketSize),
	}
}

func (b *LeakyBucket) AddTask(task Task) bool {
	// 如果木桶已经满了,返回 false
	select {
	case b.bucket <- task:
	default:
		fmt.Printf("request[id=%d] is refused\n", task.taskID)
		return false
	}
	// 如果成功入桶,调用者会等待 worker 执行结果
	fmt.Printf("request[id=%d] is run ok, resp[%v]\n", task.taskID, <-task.resultChan)
	return true
}

func (b *LeakyBucket) Start(ctx context.Context) {
	// 开启 worker 从木桶拉取任务执行
	for i := 0; i < b.WorkerNum; i++ {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					task := <-b.bucket
					task.resultChan <- task.handler()
				}
			}
		}(ctx)
	}
}

func main() {
	bucket := NewLeakyBucket(10, 4)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bucket.Start(ctx) // 开启消费者

	// 模拟20个并发请求
	var wg sync.WaitGroup
	count := 30
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(id int) {
			defer wg.Done()
			task := NewTask(id, func() Result {
				time.Sleep(300 * time.Millisecond)
				return Result{}
			})
			bucket.AddTask(task)
		}(i)
	}
	wg.Wait()
}
