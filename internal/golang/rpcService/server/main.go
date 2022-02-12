package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"study/internal/golang/rpcService"
)

func main() {
	if err := rpc.Register(rpcService.DemoService{}); err != nil {
		panic(err)
	}

	host := ":1234"

	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	log.Printf("listening port:%v", host)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
		}
		go jsonrpc.ServeConn(conn)
	}
}
