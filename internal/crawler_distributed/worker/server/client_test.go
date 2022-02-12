package main

import (
	"fmt"
	"testing"
	"time"

	"study/pkg/crawler_distributed/config"
	"study/pkg/crawler_distributed/rpcsupport"
	"study/pkg/crawler_distributed/worker"
)

func TestCrawService(t *testing.T) {
	client, err := rpcsupport.NewClient(":9000")
	if err != nil {
		return
	}

	time.Sleep(time.Second)

	req := worker.Request{
		Url: "http://album.zhenai.con/u/108906739",
		Parser: worker.SerializedParser{
			Name: "ParseProfile",
			Args: "安静的雪",
		},
	}
	resp := worker.ParseResult{}
	if err = client.Call(config.CrawServiceRpc, req, &resp); err != nil {
		t.Error(err)
	}
	fmt.Println(resp)

}
