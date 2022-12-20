package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	// 1. 建立连接
	client, err := rpc.Dial("tcp", ":1234")
	if err != nil {
		panic("connect error")
	}

	// 2. 调用服务
	var reply *string = new(string)
	err = client.Call("HelloService.Hello", "china", reply)
	if err != nil {
		panic("call error")
	}
	fmt.Println(*reply)

}
