package main

import (
	"errors"
	"log"
)

/*
	错误处理的综合示例
			defer + panic + recover
*/
func tryRecover() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			log.Println("[tryRecover] error: ", err)
		} else {
			panic(err)
		}
	}()

	panic(errors.New("this is a panic"))
}

func main() {
	tryRecover()
}
