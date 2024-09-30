package exercise

import (
	"fmt"
	"sync"
	"time"
)

/*
   1、打印12个字符串的数组,要求开启7个goroutine,并发打印的结果为 "%d goroutine: %s"
   var strSlice = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}
*/

var strSlice = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l"}

func printCharArrayConcurrent() {
	var wg sync.WaitGroup
	ch := make(chan string)

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Add(-1)
			for str := range ch {
				fmt.Printf("%d goroutine: %s,%v\n", id, str, time.Now())
			}
		}(i)
	}

	for _, v := range strSlice {
		ch <- v
	}
	close(ch)

	wg.Wait()
}

func printCharArrayConcurrentSecond() {
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Add(-1)
			// i = 1    j = 0   j = 5    j = 10
			// i = 2    j = 1   j = 6    j = 11
			// i = 3    j = 2   j = 7
			// i = 4    j = 3   j = 8
			// i = 5    j = 4   j = 9
			for j := id - 1; j < len(strSlice); j += 5 {
				fmt.Printf("%d goroutine: %s\n", id, strSlice[j])
			}
		}(i)
	}
	wg.Wait()
}
