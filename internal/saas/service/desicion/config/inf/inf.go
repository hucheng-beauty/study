package config

import "context"

// Manager 管理配置的统一入口
type Manager interface {
    Start(ctx context.Context) error // 启动所有监听和缓存机制
    Stop() error                     // 停止所有监听
    Get(key string) (any, bool)      // 获取配置缓存
}

// Source 是配置的实际来源（例如 tcc、file、memory）
type Source interface {
    // Load 用于从该配置源中加载配置
    Load(ctx context.Context, key string) (any, error)

    // Watch 注册配置变化监听，当配置变化时触发回调
    Watch(ctx context.Context, key string, onChange func(any)) error
}
