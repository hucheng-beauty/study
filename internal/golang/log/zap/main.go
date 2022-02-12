package main

import (
	"go.uber.org/zap"
)

/*
	编程范式:
		面向接口,函数式编程,并发编程
*/

/*
	Go语言并发编程
		采用 CSP 模型
		不需要锁和 callback
		并发编程 VS 并行计算
*/

func main() {

	// 使用 zap 日志库
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger.Info("incoming request",
		zap.String("path", "hello"),
		zap.Int("status", 1),
	)
}
