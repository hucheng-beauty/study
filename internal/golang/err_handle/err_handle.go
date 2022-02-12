package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
	defer
		资源管理
*/

/*
	err & panic
		意料之中的使用 error, eg.文件无法打开
		意料之外的使用 panic, eg.数据越界
*/

/*
	错误处理的综合示例
		defer + panic + recover

	type assertion:类型断言
		处理不同的错误

	函数式编程的应用: errWrapper
*/

type userError interface {
	error
	Message() string
}

func rtHello() string {
	return "hello"
}

func writeFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_EXCL|os.O_CREATE, 0666)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok {
			fmt.Println(pathError.Op, pathError.Path, pathError.Error())
		} else {
			panic(err)
		}
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 20; i++ {
		fmt.Fprintln(writer, rtHello())
	}

}

func main() {
	writeFile("test.txt")
}
