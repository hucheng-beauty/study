package main

import (
	"go.uber.org/zap"
)

// 使用 zap 日志库
func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("incoming request",
		zap.String("path", "hello"),
		zap.Int("status", 1),
	)
}
