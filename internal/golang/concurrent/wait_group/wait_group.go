package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Printf("goroutine %d\n", i)
		}(i)
	}
	wg.Wait()
}
