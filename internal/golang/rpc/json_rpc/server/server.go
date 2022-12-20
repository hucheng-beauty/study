package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct{}

func (s *HelloService) Hello(req string, reply *string) error {
	*reply = "hello, " + req
	return nil
}

func main() {
	// 1. 实例化一个 server
	listener, _ := net.Listen("tcp", ":1234")

	// 2. 注册处理逻辑 handler
	_ = rpc.RegisterName("HelloService", &HelloService{})

	for { // 防止服务端退出
		// 3. 启动服务
		conn, _ := listener.Accept()
		//rpc.ServeConn(conn)                          // 序列化协议为 Gob
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn)) // 修改序列化协议为 json
	}

}
