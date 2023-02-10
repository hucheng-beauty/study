package patterns

import (
	"fmt"
	"sync"
)

// 适用 channel 数量已知
func fanInOne(c1, c2 chan string) chan string {
	out := make(chan string)
	go func() {
		for {
			out <- <-c1
		}
	}()
	go func() {
		for {
			out <- <-c2
		}
	}()
	return out
}

// 适用 channel 数量未知
func fanInSecond(chs ...chan string) chan string {
	out := make(chan string)

	// 从 chs 读取的四种方式
	// first
	/*
		for _, ch := range chs {
			go func() {
				for {
					out <- <-ch
				}
			}()
	*/

	// second
	/*
		for _, ch := range chs {
			chCopy := ch
			go func() {
				for {
					out <- <-chCopy
				}
			}()
		}
	*/

	// third
	/*
		for _, ch := range chs {
			go func(in chan string) {
				for {
					out <- <-in
				}
			}(ch)
		}
	*/

	// fourth
	/*for i := 0; i < len(chs); i++ {
		go func(ii int) {
			for {
				out <- <-chs[ii]
			}
		}(i)
	}*/

	return out
}

func fanInForSelect(c1, c2 chan string) chan string {
	out := make(chan string)
	go func() {
		s := ""
		for {
			select {
			case s = <-c1:
				out <- s
			case s = <-c2:
				out <- s
			}
		}
	}()
	return out
}

// 工序 1 采购
func buy(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- fmt.Sprint("零件", i)
		}
	}()
	return out
}

// 工序 2 组装
func build(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "组装(" + c + ")"
		}
	}()
	return out
}

// pack 工序 3 打包
func pack(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c := range in {
			out <- "打包(" + c + ")"
		}
	}()
	return out
}

// FanIn 扇入函数（组件,把多个 channel 中的数据发送到一个 channel 中
func FanIn(ins ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)
	defer close(out)

	// 把一个 channel 中的数据发送到 out 中
	fanIn := func(in <-chan string) {
		defer wg.Add(-1)
		for v := range in {
			out <- v
		}
	}

	// 扇入,需要启动多个 goroutine 用于处于多个 channel 中的数据
	for _, cs := range ins {
		wg.Add(1)
		go fanIn(cs)
	}

	// 等待所有输入的数据ins处理完，再关闭输出out
	wg.Wait()
	return out
}

// FanInAndOut 扇入扇出模式
func FanInAndOut() {
	components := buy(10) // 采购10套配件

	// 3 班人同时组装 100 部手机
	phones1 := build(components)
	phones2 := build(components)
	phones3 := build(components)

	// fan out
	for p := range /*打包*/ pack(FanIn(phones1, phones2, phones3) /*汇聚三个channel成一个*/) {
		fmt.Println(p)
	}
}
