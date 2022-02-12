package factory_method

/*
	工厂方法模式:
		1.实现一个抽象的类
		2.实现具体的产品1
		3.实现具体的产品2
		4.实现抽象工厂类
*/

// 1.实现一个抽象的类

type Cache interface {
	Set(key, value string)
	Get(key string) string
}

// 2.实现具体的产品1

type Redis struct {
	m map[string]string
}

func NewRedis(m map[string]string) *Redis {
	return &Redis{m: m}
}

func (r *Redis) Set(key, value string) {
	r.m[key] = value
}

func (r *Redis) Get(key string) string {
	return r.m[key]
}

// 3.实现具体的产品2

type MemoryCache struct {
	m map[string]string
}

func NewMemoryCache(m map[string]string) *MemoryCache {
	return &MemoryCache{m: m}
}

func (m *MemoryCache) Set(key, value string) {
	m.m[key] = value
}

func (m *MemoryCache) Get(key string) string {
	return m.m[key]
}

// 4.实现抽象工厂类

type CacheFactory interface {
	Creat() Cache
}

// 4.1实现具体Redis工厂

type RedisCacheFactory struct{}

func (rf *RedisCacheFactory) Create() Cache {
	return &Redis{
		m: make(map[string]string),
	}
}

// 4.2实现具体Memory工厂

type MemoryCacheFactory struct{}

func (mf *MemoryCacheFactory) Create() Cache {
	return &MemoryCache{
		m: make(map[string]string),
	}
}
