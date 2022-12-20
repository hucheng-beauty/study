package main

import "github.com/gin-gonic/gin"

func pong(c *gin.Context) {
	c.JSON(200, gin.H{
		"message`": "pong",
	})
}

func main() {
	r := gin.Default()
	r.GET("/ping", pong)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
