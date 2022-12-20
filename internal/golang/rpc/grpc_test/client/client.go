package main

import (
	"context"
	"fmt"
	"study/internal/golang/rpc/grpc_test/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8088", grpc.WithInsecure())
	if err != nil {
		panic("dial error")
	}
	defer conn.Close()

	c := proto.NewGreeterClient(conn)
	helloReply, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "china"})
	if err != nil {
		panic(err)
	}
	fmt.Println(helloReply.Message)
}
