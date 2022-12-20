package factory_mode

// Cache 抽象的产品
type Cache interface {
	Set(key, value string)
	Get(key string) string
}

// CacheFactory 抽象的产品工厂
type CacheFactory interface {
	Create() (Cache, error)
}
