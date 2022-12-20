package main

import (
	"fmt"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func main() {
	// 1. 建立连接
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic("connect error")
	}

	// 2. 调用服务
	var reply *string = new(string)

	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn)) // 序列化方式为 json
	err = client.Call("HelloService.Hello", "china", reply)
	if err != nil {
		panic("call error")
	}
	fmt.Println(*reply)

}
