package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	r := gin.Default()

	// 使用 zap 日志库
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	// 增加 middleware(中间件)
	r.Use(func(c *gin.Context) {
		s := time.Now()

		// path, log latency, response code
		c.Next()
		logger.Info("incoming request",
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("elapsed", time.Now().Sub(s)),
		)

	}, func(c *gin.Context) {
		c.Set("requestId", rand.Int())

		c.Next()
	})

	// ping 请求
	r.GET("/ping", func(c *gin.Context) {
		h := gin.H{
			"message": "pong",
		}

		rid, exists := c.Get("requestId")
		if exists {
			h["requestId"] = rid
		}

		c.JSON(http.StatusOK, h)
	})

	// hello 请求
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello")
	})

	// 启动
	r.Run()
}
