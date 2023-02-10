package exercise

import (
	"fmt"
	"sync"
)

/*
   1、打印12个字符串的数组，要求开启7个goroutine，并发打印的结果为 "%d goroutine: %s"
   var strSlice = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
*/

var strSlice = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

func solutionTwo() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			// i = 1    j = 0   j = 5    j = 10
			// i = 2    j = 1   j = 6    j = 11
			// i = 3    j = 2   j = 7
			// i = 4    j = 3   j = 8
			// i = 5    j = 4   j = 9
			for j := id - 1; j < len(strSlice); j += 5 {
				fmt.Printf("%d goroutine: %s\n", id, strSlice[j])
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func solutionOne() {
	var wg sync.WaitGroup
	ch := make(chan string)

	// read
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Add(-1)
			for str := range ch {
				fmt.Printf("%d goroutine: %s\n", id, str)
			}
		}(i)
	}

	// write
	for _, v := range strSlice {
		ch <- v
	}
	close(ch)

	wg.Wait()
}

func asStream(done <-chan struct{}, values ...any) <-chan any {
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

func main() {
	// done := make(chan struct{})
	// for s := range asStream(done, strSlice) {
	//     fmt.Println(s)
	// }
	solutionOne()
}
