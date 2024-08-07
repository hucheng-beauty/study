package exercise

import (
	"fmt"
	"sync"
)

/*
	编写一个可以从 1 到 n 输出代表这个数字的字符串的程序,要求:
		如果这个数字可以被 3 整除,输出 "fizz"。
		如果这个数字可以被 5 整除,输出 "buzz"。
		如果这个数字可以同时被 3 和 5 整除,输出 "fizzbuzz"。
	例如,当 n = 15,输出: 1, 2, fizz, 4, buzz, fizz, 7, 8, fizz, buzz, 11, fizz, 13, 14, fizzbuzz。
		假设有这么一个结构体:
		type FizzBuzz struct {}
		func (fb *FizzBuzz) fizz() {}
		func (fb *FizzBuzz) buzz() {}
		func (fb *FizzBuzz) fizzbuzz() {}
		func (fb *FizzBuzz) number() {}
	请你实现一个有四个线程的多协程版 FizzBuzz,同一个 FizzBuzz 对象会被如下四个协程使用:
	协程 A 将调用 fizz() 来判断是否能被 3 整除,如果可以,则输出 fizz
	协程 B 将调用 buzz() 来判断是否能被 5 整除,如果可以,则输出 buzz
	协程 C 将调用 fizzbuzz() 来判断是否同时能被 3 和 5 整除,如果可以,则输出 fizzbuzz
	协程 D 将调用 number() 来实现输出既不能被 3 整除也不能被 5 整除的数字
*/

type FizzBuzz struct {
	n  int
	ch chan int
	wg sync.WaitGroup
}

func New(n int) *FizzBuzz {
	return &FizzBuzz{
		n:  n,
		ch: make(chan int, 1),
	}
}

func (fb *FizzBuzz) start() {
	fb.wg.Add(4)
	go fb.fizz()
	go fb.buzz()
	go fb.fizzbuzz()
	go fb.number()
	fb.ch <- 1
	fb.wg.Wait()
}

func (fb *FizzBuzz) fizz() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%3 == 0 {
			if v%5 == 0 {
				fb.ch <- v
				continue
			}
			if v == fb.n {
				fmt.Print(" fizz。")
			} else {
				fmt.Print(" fizz,")
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) buzz() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%5 == 0 {
			if v%3 == 0 {
				fb.ch <- v
				continue
			}
			if v == fb.n {
				fmt.Print(" buzz。")
			} else {
				fmt.Print(" buzz,")
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) fizzbuzz() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%5 == 0 && v%3 == 0 {
			if v == fb.n {
				fmt.Print(" fizzbuzz。")
			} else {
				fmt.Print(" fizzbuzz,")
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}

func (fb *FizzBuzz) number() {
	defer fb.wg.Done()
	for v := range fb.ch {
		if v > fb.n {
			fb.ch <- v
			return
		}
		if v%5 != 0 && v%3 != 0 {
			if v == fb.n {
				fmt.Printf(" %d。", v)
			} else {
				fmt.Printf(" %d,", v)
			}
			fb.ch <- v + 1
			continue
		}
		fb.ch <- v
	}
}
