package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"study/internal/golang/rpc/stream_grpc/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	client := proto.NewGreeterClient(conn)

	// 服务端流模式
	// resp, err := client.GetStream(context.Background(), &proto.StreamReqData{Data: "hello"})
	// if err != nil {
	//	return
	// }
	// for {
	//	data, err := resp.Recv()
	//	if err != nil {
	//		fmt.Println(err)
	//		break
	//	}
	//	fmt.Println(data)
	// }

	// 客户端流模式
	// putStream, _ := client.PutStream(context.Background())
	// i := 0
	// for {
	//	i++
	//	_ = putStream.Send(&proto.StreamReqData{Data: fmt.Sprintf("%v", time.Now().Unix())})
	//	time.Sleep(time.Second)
	//
	//	if i > 10 {
	//		break
	//	}
	// }

	// 双向流模式
	allStream, _ := client.AllStream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		for {
			// receive
			data, _ := allStream.Recv()
			fmt.Println("接受服务端消息", data)
			time.Sleep(time.Millisecond * 100)
		}
	}()
	go func() {
		for {
			// send
			_ = allStream.Send(&proto.StreamReqData{Data: "hello server"})
			time.Sleep(time.Millisecond * 100)
		}
	}()
	wg.Wait()
}
