package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {

	arr1 := []int{1, 3, 5}
	arr2 := []int{2, 4, 6}
	wg.Add(2)

	ch1 := make(chan struct{})
	ch2 := make(chan struct{}, 1)

	ch2 <- struct{}{}
	go func() {
		defer wg.Done()
		for _, v := range arr1 {
			<-ch2
			print(v)
			ch1 <- struct{}{}
		}

	}()
	go func() {
		defer wg.Done()
		for _, v := range arr2 {
			<-ch1
			print(v)
			ch2 <- struct{}{}
		}
	}()

	wg.Wait()
}
