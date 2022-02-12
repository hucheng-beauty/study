package factory_method

import (
	"fmt"
	"testing"
)

func TestRedis(t *testing.T) {
	redis := NewRedis(map[string]string{})

	redis.Set("name", "redis_cache")
	fmt.Println(redis.Get("name"))
}

func TestMemoryCache(t *testing.T) {
	memoryCache := NewMemoryCache(map[string]string{})

	memoryCache.Set("name", "memory_cache")
	fmt.Println(memoryCache.Get("name"))
}

func TestRedisCacheFactory(t *testing.T) {
	redisCacheFactory := RedisCacheFactory{}
	redisCache := redisCacheFactory.Create()

	redisCache.Set("name", "redis_factory")
	fmt.Println(redisCache.Get("name"))
}

func TestMemoryCacheFactory(t *testing.T) {
	memoryCacheFactory := MemoryCacheFactory{}
	memoryCache := memoryCacheFactory.Create()

	memoryCache.Set("name", "memory_factory")
	fmt.Println(memoryCache.Get("name"))
}
