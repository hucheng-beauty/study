package abstract_factory

/*
	工厂方法模式:
		抽象创建工厂	<-----(实现)         具体产品1工厂
										   ｜
										   ｜(创建)
	  	抽象的产品	<-----(实现) 		具体的产品1


		优点: 保持了简单工厂模式的优点,克服了它的缺点
		缺点: 在添加新产品时,在一定程度上增加了系统的复杂性
		适合: 客户端不需要知道具体产品类的类名,只需知道所对应的工厂即可

*/

// Cache 1.实现一个抽象的类
type Cache interface {
	Set(key, value string)
	Get(key string) string
}

// RedisCache 2.实现具体的产品1
type RedisCache struct {
	m map[string]string
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		m: make(map[string]string),
	}
}

func (rc *RedisCache) Set(key, value string) {
	rc.m[key] = value
}

func (rc *RedisCache) Get(key string) string {
	return rc.m[key]
}

// MemoryCache 3.实现具体的产品2
type MemoryCache struct {
	m map[string]string
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		m: make(map[string]string),
	}
}

func (mc *MemoryCache) Set(key, value string) {
	mc.m[key] = value
}

func (mc *MemoryCache) Get(key string) string {
	return mc.m[key]
}

// CacheFactory 4.实现工厂类
type CacheFactory interface {
	Create() (Cache, error)
}

type RedisCacheFactory struct{}

func NewRedisCacheFactory() *RedisCacheFactory {
	return &RedisCacheFactory{}
}

func (rcf *RedisCacheFactory) Create() (Cache, error) {
	return &RedisCache{
		m: make(map[string]string),
	}, nil
}

type MemoryCacheFactory struct{}

func NewMemoryCacheFactory() *MemoryCacheFactory {
	return &MemoryCacheFactory{}
}

func (mcf *MemoryCacheFactory) Create() (Cache, error) {
	return &MemoryCache{
		m: make(map[string]string),
	}, nil
}
