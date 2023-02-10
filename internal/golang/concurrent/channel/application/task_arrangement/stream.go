package task_arrangement

// 把 Channel 当作流式管道使用的方式,也就是把 Channel 看作流(Stream),提供跳过几个元素,或者是只取其中的几个元素等方法
func asStream(done <-chan struct{}, values ...interface{}) <-chan interface{} {
	out := make(chan interface{}) // 创建一个 unbuffered 的 channel
	go func() {
		// 启动一个 goroutine，往 out 中塞数据
		defer close(out) // 退出时关闭 chan

		for _, v := range values { // 遍历数组
			select {
			case <-done:
				return
			case out <- v: // 将数组元素塞入到 chan 中
			}
		}
	}()
	return out
}

func takeN(done <-chan struct{}, inStream <-chan interface{}, num int) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)

		for i := 0; i < num; i++ { // 只读取前 num 个元素
			select {
			case <-done:
				return
			case outStream <- <-inStream: // 从输入流中读取元素
			}
		}
	}()
	return outStream
}
