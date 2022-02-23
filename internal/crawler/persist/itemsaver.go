package persist

import (
	"errors"
	"log"

	"study/internal/crawler/engine"

	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
)

// ItemSaver :Send data to ElasticSearch
func ItemSaver(index string) (chan engine.Item, error) {
	// Must turn off sniff in docker
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			itemCount++
			// storage data to ElasticSearch
			if err = Save(client, index, item); err != nil {
				log.Printf("Item Saver error saving item %v: %v", item, err)
			}
		}
	}()
	return out, nil
}

func Save(client *elastic.Client, index string, item engine.Item) error {
	if item.Type == "" {
		return errors.New("must supply Type")
	}

	itemService := client.Index().
		Index(index).
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
