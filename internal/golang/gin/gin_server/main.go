package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	AlarmCount int64
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

	})

	r.POST("/", func(c *gin.Context) {
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("[Handler] ReadAll error:", err)
			c.String(http.StatusBadRequest, err.Error())
		}
		atomic.AddInt64(&AlarmCount, 1)

		c.JSON(http.StatusOK, gin.H{
			"message": "ok!",
			"request": string(bodyBytes),
		})
	})

	r.GET("/alarm_count", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":     "ok!",
			"alarm_count": AlarmCount,
		})
	})

	// 启动
	r.Run(":7888")
}
