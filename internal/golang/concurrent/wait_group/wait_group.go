package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}
	workerCount := 3
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Add(-1)
			fmt.Printf("goroutine %d\n", id)
		}(i)
	}
	wg.Wait()
}
