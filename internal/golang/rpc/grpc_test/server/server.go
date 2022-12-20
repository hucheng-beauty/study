package main

import (
	"context"
	"net"
	"study/internal/golang/rpc/grpc_test/proto"

	"google.golang.org/grpc"
)

type Server struct{}

func (s Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{Message: "hello " + request.Name}, nil
}

func main() {
	g := grpc.NewServer()
	proto.RegisterGreeterServer(g, &Server{})

	listener, err := net.Listen("tcp", "0.0.0.0:8088")
	if err != nil {
		panic("listen error")
	}

	_ = g.Serve(listener)
}
