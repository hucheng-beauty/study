package internal

import (
    "context"
    "log"
    "sync"
    "time"

    inf "study/internal/saas/service/desicion/config/inf"
)

type ConfigManagerImpl struct {
    sources []inf.Source
    cache   sync.Map
    cancel  context.CancelFunc
}

func NewConfigManagerImpl(sources ...inf.Source) *ConfigManagerImpl {
    return &ConfigManagerImpl{sources: sources}
}

// Start 启动配置加载和监听
func (m *ConfigManagerImpl) Start(ctx context.Context) error {
    ctx, cancel := context.WithCancel(ctx)
    m.cancel = cancel

    for _, src := range m.sources {
        go m.start(ctx, src)
    }

    return nil
}

func (m *ConfigManagerImpl) start(ctx context.Context, src inf.Source) {
    keys := m.getKeysToWatch()

    // 1️⃣ 初始化加载
    for _, key := range keys {
        val, err := src.Load(ctx, key)
        if err != nil {
            log.Printf("[ConfigManager] Load %s failed: %v", key, err)
            continue
        }
        m.cache.Store(key, val)
        log.Printf("[ConfigManager] 初始加载成功: %s", key)
    }

    // 2️⃣ 注册 Watch
    for _, key := range keys {
        key := key
        err := src.Watch(ctx, key, func(v any) {
            m.cache.Store(key, v)
            log.Printf("[ConfigManager] 配置更新: %s", key)
        })
        if err != nil {
            log.Printf("[ConfigManager] Watch %s failed: %v", key, err)
        }
    }

    // 3️⃣ 可选：周期性校验（保证长期一致性）
    go func() {
        ticker := time.NewTicker(10 * time.Minute)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                for _, key := range keys {
                    val, err := src.Load(ctx, key)
                    if err == nil {
                        m.cache.Store(key, val)
                    }
                }
            }
        }
    }()
}

func (m *ConfigManagerImpl) Stop() error {
    if m.cancel != nil {
        m.cancel()
    }
    return nil
}

func (m *ConfigManagerImpl) Get(key string) (any, bool) {
    return m.cache.Load(key)
}

// getKeysToWatch 模拟从初始化文件或配置中读取要监听的配置目录
func (m *ConfigManagerImpl) getKeysToWatch() []string {
    // 实际可从本地 YAML 配置读取
    return []string{
        "speech/configs/platform",
        "speech/configs/service",
    }
}
