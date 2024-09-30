package crawler

import (
	"context"
	"errors"
	"log"

	"github.com/olivere/elastic/v7"
)

// ItemSaver send data to ElasticSearch
func ItemSaver(index string) (chan Item, error) {
	// Must turn off sniff in docker
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	out := make(chan Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			itemCount++
			// Storage data to elastic_search
			if err = save(client, index, item); err != nil {
				log.Printf("saving item %v, error:%v", item, err)
			}
		}
	}()
	return out, nil
}

func save(client *elastic.Client, index string, item Item) error {
	if item.Type == "" {
		return errors.New("must supply type")
	}

	itemService := client.Index().Index(index).
		OpType(item.Type).
		Id(item.Id).
		BodyJson(item)
	if item.Id != "" {
		itemService.Id(item.Id)
	}

	_, err := itemService.Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
