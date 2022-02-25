package patterns

import "reflect"

// 扇入模式
// 在软件工程中，模块的扇入是指有多少个上级模块调用它。
// 而对于这里的扇入模式来说，就是指有多个源 Channel 输入、一个目的 Channel 输出的情况。
// 扇入比就是源 Channel 数量比 1。每个源 Channel 的元素都会发送给目标 Channel，
// 相当于目标 Channel 的 receiver 只需要监听目标 Channel，就可以接收所有发送给源 Channel 的数据

func fanInReflect(chs ...chan interface{}) chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		// 构造 cases
		var cases []reflect.SelectCase
		for ch := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			})
		}

		// 选一个可用的
		for len(cases) > 0 {
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases, cases[chosen+1:]...)
			}
			out <- value.Interface()
		}
	}()
	return out
}

func fanInRecursion(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		c := make(chan interface{})
		close(c)
		return c
	case 1:
		return chs[0]
	case 2:
		return merge(chs[0], chs[1])
	default:
		m := len(chs) / 2
		return merge(fanInRecursion(chs[:m]...), fanInRecursion(chs[m:]...))
	}
}

func merge(a, b <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		for a != nil || b != nil { //只要还有可读的 chan
			select {
			case v, ok := <-a:
				if !ok { // a 已关闭，设置为 nil
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok { // b 已关闭，设置为 nil
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
