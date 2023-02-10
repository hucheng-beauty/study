package task_arrangement

// 扇出模式只有一个输入源 Channel,有多个目标 Channel，
// 扇出比就是 1 比目标 Channel 数的值,经常用在设计模式中的
// 观察者模式中(观察者设计模式定义了对象间的一种一对多的组合关系;一个对象的状态发生变化时,所有依赖于它的对象都会得到通知并自动刷新);
// 在观察者模式中,数据变动后,多个观察者都会收到这个变更信号。
// fanOut 扇出
func fanOut(ch <-chan interface{}, outs []chan interface{}, async bool) {
	go func() {
		defer func() { // 退出时关闭所有的输出chan
			for i := 0; i < len(outs); i++ {
				close(outs[i])
			}
		}()

		for value := range ch { // 从输入chan中读取数据
			v := value
			for i := 0; i < len(outs); i++ {
				i := i
				if async { // 异步
					go func() {
						outs[i] <- v // 放入到输出 chan 中,异步方式
					}()
				} else {
					outs[i] <- v // 放入到输出 chan 中，同步方式
				}
			}
		}
	}()
}
