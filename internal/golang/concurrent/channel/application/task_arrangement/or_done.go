package task_arrangement

import "reflect"

// channel 的任务编排模式

// Or-Done 模式
// 我们会使用“信号通知”实现某个任务执行完成后的通知机制，
// 在实现时，我们为这个任务定义一个类型为 chan struct{}类型的 done 变量，
// 等任务结束后，我们就可以 close 这个变量，然后，其它 receiver 就会收到这个通知。
// 这是有一个任务的情况，如果有多个任务，只要有任意一个任务执行完，我们就想获得这个信号，
// 这就是 Or-Done 模式。
// 递归实现等待信号
func orRecursion(channels ...<-chan interface{}) <-chan interface{} {
	// 特殊情况，只有零个或者1个chan
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2: // 2个也是一种特殊情况
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: // 超过两个，二分法递归处理
			m := len(channels) / 2
			select {
			case <-orRecursion(channels[:m]...):
			case <-orRecursion(channels[m:]...):
			}
		}
	}()

	return orDone
}

// 反射实现等待信号
func orReflect(channels ...<-chan interface{}) <-chan interface{} {
	//特殊情况，只有0个或者1个
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)
		// 利用反射构建SelectCase
		var cases []reflect.SelectCase
		for _, c := range channels {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		// 随机选择一个可用的case
		reflect.Select(cases)
	}()

	return orDone
}
