package main

import (
	"fmt"
	"time"
)

type Token struct{}

func main() {
	chs := make([]chan Token, 4)

	for i := 0; i < 4; i++ {
		chs[i] = make(chan Token)
		go func(id int) {
			for {
				token := <-chs[id]
				fmt.Println(id + 1)
				time.Sleep(100 * time.Millisecond)
				chs[(id+1)%4] <- token
			}
		}(i)
	}

	// 首先把令牌交给第一个worker
	// chs[0] ==> go chs[1] ==> go chs[2] ==> go chs[3] ==> go chs[0]
	chs[0] <- struct{}{}

	select {}
}
