package client

import (
	"log"

	"study/internal/crawler/engine"
	"study/internal/crawler_distributed/config"
	"study/internal/crawler_distributed/rpcsupport"
)

// ItemSaver :Send data to ElasticSearch
func ItemSaver(host string) (chan engine.Item, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			itemCount++

			// Call RPC to storage data to ElasticSearch
			result := ""
			if err := client.Call(config.ItemSaverRpc, item, result); err != nil {
				log.Printf("Item Saver error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}
