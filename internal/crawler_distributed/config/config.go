package config

const (
	// Service ports
	ItemSaverPort = 1234
	WorkerPort0   = 9000

	// ElasticSearch
	ElasticIndex = "dating_profile"

	// RPC Endpoints
	ItemSaverRpc   = "ItemSaverService.Save"
	CrawServiceRpc = "CrawlService.Process"
)
