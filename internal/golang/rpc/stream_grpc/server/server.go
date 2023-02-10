package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	"study/internal/golang/rpc/stream_grpc/proto"

	"google.golang.org/grpc"
)

const PORT = ":50052"

type server struct{}

func (s *server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		_ = res.Send(&proto.StreamRespData{
			Data: fmt.Sprintf("%v", time.Now().Unix()),
		})
		time.Sleep(1 * time.Second)

		if i > 10 {
			break
		}
	}
	return nil
}

func (s *server) PutStream(clientStream proto.Greeter_PutStreamServer) error {
	for {
		data, err := clientStream.Recv()
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(data)
		}
	}

	return nil
}

func (s *server) AllStream(allStream proto.Greeter_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		// send
		for {
			_ = allStream.Send(&proto.StreamRespData{Data: "hello client"})
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		for true {
			// receive
			data, _ := allStream.Recv()
			fmt.Println("接受客户端消息", data)
			time.Sleep(time.Millisecond * 100)
		}
	}()

	wg.Wait()
	return nil
}

func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		return
	}

	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	s.Serve(listen)
}
