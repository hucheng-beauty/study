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
    "Redis":  NewRedis(),
    "Memory": NewMemory(),
}

func SelectCache(cacheType string) Cache {
    // SelectCache 适合对象中无字段
    cache, ok := CachesFactory[cacheType]
    if !ok {
        return CachesFactory["Redis"] // 默认使用 Redis
    }
    return cache
}

func NewCache(cacheType string) (Cache, error) {
    // NewCache 适合对象中有字段
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

func NewRedisFactory() *RedisFactory {
    return &RedisFactory{}
}

func (r *RedisFactory) Create() (Cache, error) {
    return NewRedis(), nil
}

type MemoryFactory struct{}

func NewMemoryFactory() *MemoryFactory {
    return &MemoryFactory{}
}

func (m *MemoryFactory) Create() (Cache, error) {
    return NewMemory(), nil
}
