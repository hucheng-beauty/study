package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"

	"study/internal/golang/rpcService"
)

func main() {
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	client := jsonrpc.NewClient(conn)

	var result float64
	if err = client.Call("DemoService.Div", rpcService.Args{A: 2, B: 4}, &result); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	if err = client.Call("DemoService.Div", rpcService.Args{A: 2, B: 0}, &result); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}
}
