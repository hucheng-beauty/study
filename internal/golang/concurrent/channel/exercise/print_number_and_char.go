package exercise

import (
	"fmt"
	"time"
)

func PrintNumberAndChar() {
	ch1 := make(chan struct{})
	ch2 := make(chan struct{}, 1)

	ch2 <- struct{}{}
	go func() {
		for i := 1; i <= 26; i++ {
			<-ch2
			fmt.Printf("%d\n", i)
			ch1 <- struct{}{}
		}
	}()

	go func() {
		for i := 'a'; i <= 'z'; i++ {
			<-ch1
			fmt.Printf("%s\n", string(i))
			ch2 <- struct{}{}
		}
	}()

	time.Sleep(10 * time.Second)
}

func printNumberAndChar(id int, ch chan struct{}, nextChan chan struct{}) {
	if id%2 == 0 {
		for {
			for i := 1; i <= 26; i++ {
				<-ch
				fmt.Printf("%d\n", i)
				time.Sleep(100 * time.Millisecond)
				nextChan <- struct{}{}
			}
		}
	} else {
		for {
			for i := 'a'; i <= 'z'; i++ {
				<-ch
				fmt.Printf("%s\n", string(i))
				nextChan <- struct{}{}
			}
		}
	}
}

func PrintNumberAndCharSecond() {
	chs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
	}

	for i := 0; i < 2; i++ {
		go printNumberAndChar(i, chs[i], chs[(i+1)%2])
	}

	chs[0] <- struct{}{}
	time.Sleep(10 * time.Second)
}
