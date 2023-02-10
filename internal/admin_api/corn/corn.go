package corn

import (
	"time"

	"github.com/go-co-op/gocron"
)

func LocalCache() {
	s := gocron.NewScheduler(time.Local)
	s.Every(cacheInterval).Second().StartImmediately().Do(localCache)
	s.StartAsync()
}

func localCache() {
	SynUserCache()
}
