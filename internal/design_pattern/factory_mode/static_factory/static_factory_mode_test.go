package static_factory

import (
	"fmt"
)

func main() {
	cacheType := "Redis"
	cacheFactory := CacheFactory{}
	cache, err := cacheFactory.Create(cacheType)
	if err != nil {
		return
	}
	cache.Set("Id", "001")
	fmt.Println(cache.Get("Id"))

}
