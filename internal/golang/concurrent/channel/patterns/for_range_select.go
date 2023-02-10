package patterns

import (
	"fmt"
	"time"
)

// ForRangeSelect for range select 有限循环模式
func ForRangeSelect() {
	quit := make(chan struct{})
	for _, i := range []int{1, 2, 3} {
		fmt.Println(i)
		go func(i int) {
			// send
			if i == 3 {
				quit <- struct{}{}
			}
		}(i)
	}

	time.Sleep(time.Millisecond)

	for {
		select {
		// receive
		case <-quit:
			fmt.Println("over")
			return
		}
	}
}
