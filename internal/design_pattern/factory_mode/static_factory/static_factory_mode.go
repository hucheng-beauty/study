package static_factory

import "errors"

/*
	简单工厂模式(静态工厂方法模式):
		简单工厂模式专门定义一个类来负责创建其他类的实例,
		被创建的实例通常都具有共同的父类,可以根据参数的不同返回不同类的实例

	简单工厂模式:
		1.实现一个抽象的类
		2.实现具体的产品1
		3.实现具ç体的产品2
		4.实现简单工厂类

		优点: 实现了解耦
		缺点: 违背了开闭原则
		适合:
*/

// Cache 1.实现一个抽象的类
type Cache interface {
	Set(key, value string)
	Get(key string) string
}

// Redis 2.实现具体的产品Redis
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

// MemoryCache 3.实现具体的产品 Memory
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

// CachesFactory 4.1实现简单工厂类:适合对象中无字段
var CachesFactory = map[string]Cache{
	"Redis":       new(Redis),
	"MemoryCache": new(MemoryCache),
}

func SelectCache(cacheType string) Cache {
	cache, ok := CachesFactory[cacheType]
	if !ok {
		return CachesFactory["Redis"]

	}

	return cache
}

// CacheFactory 4.2实现简单工厂类
type CacheFactory struct{}

func (CacheFactory) Create(cacheType string) (Cache, error) {
	switch cacheType {
	case "Redis":
		return NewRedis(make(map[string]string)), nil
	case "MemoryCache":
		return NewMemoryCache(make(map[string]string)), nil
	default:
		return nil, errors.New("cache type invalid")
	}
}
