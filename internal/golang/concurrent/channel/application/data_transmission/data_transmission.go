package main

import (
	"fmt"
	"time"
)

type Token struct{}

func main() {
	chs := []chan Token{
		make(chan Token),
		make(chan Token),
		make(chan Token),
		make(chan Token),
	}
	/*
		go chs[0] ==> go chs[1] ==> go chs[2] ==> go chs[3] ==> go chs[0]
	*/

	for i := 0; i < 4; i++ {
		go func(id int) {
			for {
				token := <-chs[id]
				fmt.Println(id + 1)
				time.Sleep(100 * time.Millisecond)
				chs[(id+1)%4] <- token
			}
		}(i)
	}

	//首先把令牌交给第一个worker
	chs[0] <- struct{}{}

	select {}
}
