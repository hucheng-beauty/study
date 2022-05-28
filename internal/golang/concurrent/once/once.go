package main

import (
	"sync"
	"sync/atomic"
)

/*
	Once 常用来初始化单例资源,或者并发访问只需初始化一次的共享资源,或者在测试的时候初始化一次测试资源

	常见的 Once 错误
		死锁: 在 Do 方法中 调用自身
		未初始化: Do 执行完的初始化任务执行失败或者 panic,
			    此时, Once 还是会认为初次执行已经成功了,即使再次调用 Do 方法,也不会再次执行 Do 方法
			解决方案: 自己定制化开发类似 Once 的功能,Do传入的 func 返回 error,一直初始化至成功为止
*/

// Once 使用双检查机制实现 Once
type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	// 双检查
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

func main() {
	var once sync.Once
	once.Do(func() {
		println("hello world! => 1")
	})

	once.Do(func() {
		println("hello world! => 2") // 无输出
	})

	var o Once
	o.Do(func() {
		println("hello world! => 1")
	})

	o.Do(func() {
		println("hello world! => 2") // 无输出
	})
}
