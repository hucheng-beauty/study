package main

import (
	"flag"
	"fmt"
	"log"

	"study/pkg/crawler_distributed/config"
	"study/pkg/crawler_distributed/persist"
	"study/pkg/crawler_distributed/rpcsupport"

	"github.com/olivere/elastic/v7"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}

	host := fmt.Sprintf(":%d", config.ItemSaverPort)
	log.Fatal(serveRpc(host, config.ElasticIndex))

	/*
		if err := serveRpc(":1234", "dating_profile"); err != nil {
			panic(err)
		}
	*/
}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}

	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
