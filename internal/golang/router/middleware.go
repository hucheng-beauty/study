package router

import (
    "fmt"

    "github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("[Auth] 检查内部权限")
        c.Next()
    }
}

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        fmt.Println("[Logger] 记录请求日志")
        c.Next()
    }
}
