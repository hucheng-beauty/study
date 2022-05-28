package main

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// DoneOnce 是一个扩展的 sync.Once 类型,提供了一个Done方法
type DoneOnce struct {
	sync.Once
}

// Done 返回此Once是否执行过
// 如果执行过则返回true
// 如果没有执行过或者正在执行,返回false
func (o *DoneOnce) Done() bool {
	return atomic.LoadUint32((*uint32)(unsafe.Pointer(&o.Once))) == 1
}

/*func main() {
	var flag DoneOnce
	fmt.Println(flag.Done()) //false

	flag.Do(func() {
		time.Sleep(time.Second)
	})

	fmt.Println(flag.Done()) //true
}*/
