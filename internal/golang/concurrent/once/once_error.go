package main

import (
	"sync"
	"sync/atomic"
)

// ErrorOnce 一个功能更加强大的Once
type ErrorOnce struct {
	m    sync.Mutex
	done uint32
}

// Do 方法会把这个error返回给调用者;传入的函数f有返回值error,如果初始化失败,需要返回失败的error
func (o *ErrorOnce) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 { //fast path
		return nil
	}
	return o.slowDo(f)
}

// 如果还没有初始化
func (o *ErrorOnce) slowDo(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 { // 双检查，还没有初始化
		err = f()
		if err == nil { // 初始化成功才将标记置为已初始化
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}
