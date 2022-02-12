package main

import (
	"fmt"
	"time"

	"study/internal/golang/interface/mock"
	"study/internal/golang/interface/real"
)

/*
	接口:
		接口由使用者定义

	接口的实现:
		接口的实现是隐式的
		只要实现接口里的方法

	接口变量里面有什么?
		实现者的类型
		实现者的值(实现者的指针==>实现者)

		接口变量自带指针
		接口变量同样采用值传递,几乎不需要使用接口的指针
		指针接收者实现只能以指针方式使用;值接收者都可以

	查看接口变量
		表示任何变量 interface{}
		Type Assertion 接口断言
		Type Switch

	接口组合(由使用者组合)

	常用的系统接口
		Stringer
		Reader/Writer
*/

const url = "https://www.imooc.com"

type Retriever interface {
	Get(url string) string
}

func download(r Retriever) string {
	return r.Get(url)
}

func inspect(r Retriever) {
	fmt.Printf("%T %v\n", r, r)

	// Type Switch
	switch v := r.(type) {
	case mock.Retriever:
		fmt.Println("Contents:", v.Content)
	case *real.Retriever:
		fmt.Println("UserAgent", v.UserAgent)
	}
}

func main() {
	var r Retriever
	r = mock.Retriever{"This is a fake website."}
	r = &real.Retriever{
		UserAgent: "Mozilla/5.0",
		TimeOut:   time.Minute,
	}
	inspect(r)

	// Type assertion,接口断言
	if mockRetriever, ok := r.(mock.Retriever); ok {
		fmt.Println(mockRetriever.Content)
	} else {
		fmt.Println("Not a mock retriever.")
	}

	// fmt.Println(download(r))
}
