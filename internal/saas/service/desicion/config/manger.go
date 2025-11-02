package config

import (
    "time"

    "study/internal/saas/config"
    inf "study/internal/saas/service/desicion/config/inf"
    "study/internal/saas/service/desicion/config/internal"
)

// NewManager 创建 ConfigManager 实例
func NewManager(sources ...inf.Source) ConfigManager {
    return internal.NewConfigManagerImpl(sources...)
}
