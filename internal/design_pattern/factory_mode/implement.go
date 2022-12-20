package factory_mode

/*
	简单工厂
		具体的产品(Redis、Memory) + 产品工厂(SelectCache/CreateCache)
	工厂方法
		产品接口(Cache) + 具体的产品(Redis、Memory) + 具体的产品工厂(NewRedisFactory、NewMemoryFactory)
	抽象工厂
		产品接口(Cache) + 具体的产品(Redis、Memory) + 产品工厂接口(CacheFactory)  + 具体的产品工厂(RedisFactory、MemoryFactory)
*/

type Redis struct {
	m map[string]string
}

func (rc *Redis) Set(key, value string) {
	rc.m[key] = value
}

func (rc *Redis) Get(key string) string {
	return rc.m[key]
}

func NewRedis() *Redis {
	return &Redis{m: make(map[string]string)}
}

type Memory struct {
	m map[string]string
}

func (mc *Memory) Set(key, value string) {
	mc.m[key] = value
}

func (mc *Memory) Get(key string) string {
	return mc.m[key]
}

func NewMemory() *Memory {
	return &Memory{m: make(map[string]string)}
}

var CachesFactory = map[string]Cache{
	"Redis":  &Redis{m: make(map[string]string)},
	"Memory": &Memory{m: make(map[string]string)},
}

// SelectCache 适合对象中无字段
func SelectCache(cacheType string) Cache {
	cache, ok := CachesFactory[cacheType]
	if !ok {
		return CachesFactory["Redis"] // 默认使用 Redis
	}
	return cache
}

// NewCache 适合对象中有字段
func CreateCache(cacheType string) (Cache, error) {
	switch cacheType {
	case "Redis":
		return NewRedis(), nil
	case "Memory":
		return NewMemory(), nil
	default:
		return NewRedis(), nil
	}
}

type RedisFactory struct{}

func (rf *RedisFactory) Create() (Cache, error) {
	return NewRedis(), nil
}

func NewRedisFactory() *RedisFactory {
	return &RedisFactory{}
}

type MemoryFactory struct{}

func (mf *MemoryFactory) Create() (Cache, error) {
	return NewMemory(), nil
}

func NewMemoryFactory() *MemoryFactory {
	return &MemoryFactory{}
}
