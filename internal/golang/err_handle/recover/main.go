package main

import (
	"errors"
	"log"
)

/*
错误处理的综合示例

	defer + panic + recover
*/
func Recover() {
	defer func() {
		var err error
		if errors.As(recover().(error), &err) {
			log.Println("[recover] error: ", err)
		}
	}()

	panic(errors.New("this is a panic"))
}

func main() {
	Recover()
}
