package main

import (
	"flag"
	"log"
	"net/rpc"
	"strings"

	"study/internal/crawler/engine"
	"study/internal/crawler/scheduler"
	"study/internal/crawler/zhenai/parser"
	client1 "study/internal/crawler_distributed/persist/client"
	"study/internal/crawler_distributed/rpcsupport"
	"study/internal/crawler_distributed/worker/client"
)

const (
	url = "https://www.zhenai.com/zhenghun"
)

var (
	itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")

	workerHosts = flag.String("worker_hosts", "", "worker hosts(comma separated)")
)

func main() {
	flag.Parse()

	// config ES index
	itemChan, err := client1.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}

	pool := createClientPool(strings.Split(*workerHosts, ","))

	processor := client.CreateProcessor(pool)

	// config Engine
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      10,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}

	// start crawler
	e.Run(engine.Request{
		Url:    url,
		Parser: engine.NewFuncParser(parser.ParseCityList, "ParseCityList"),
	})
}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		c, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, c)
			log.Printf("Connecting to %v", h)
		} else {
			log.Printf("error to connecting to %s: %v", h, err)
		}
	}

	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, c := range clients {
				out <- c
			}
		}
	}()
	return out
}
