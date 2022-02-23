package client

import (
	"net/rpc"

	"study/internal/crawler/engine"
	"study/internal/crawler_distributed/config"
	"study/internal/crawler_distributed/worker"
)

func CreateProcessor(pool chan *rpc.Client) engine.Processor {
	return func(req engine.Request) (engine.ParseResult, error) {
		sReq := worker.SerializedRequest(req)
		resp := worker.ParseResult{}

		c := <-pool
		if err := c.Call(config.CrawServiceRpc, sReq, &resp); err != nil {
			return engine.ParseResult{}, err
		}
		return worker.DeserializedResult(resp), nil
	}
}
