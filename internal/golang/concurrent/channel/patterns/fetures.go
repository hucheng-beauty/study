package patterns

import (
	"fmt"
	"time"
)

// 洗菜
func wash() <-chan string {
	out := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		out <- "洗好的菜"
	}()
	return out
}

// 烧水
func water() <-chan string {
	out := make(chan string)
	go func() {
		time.Sleep(5 * time.Second)
		out <- "烧开的水"
	}()
	return out
}

// Futures 模式
func Futures() {
	fmt.Println("已经安排好洗菜和烧水了，我先开一局")
	time.Sleep(2 * time.Second)

	fmt.Println("要做火锅了，看看菜和水好了吗")

	fmt.Println("准备好了，可以做火锅了:", <-wash() /*洗菜*/, <-water() /*烧水*/)

}
