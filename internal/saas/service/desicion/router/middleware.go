package router

import (
    "context"
    "fmt"

    "github.com/cloudwego/hertz/pkg/app"
)

func Auth() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        fmt.Println("[Auth] 权限校验")
        c.Next(ctx)
    }
}

func Logger() app.HandlerFunc {
    return func(ctx context.Context, c *app.RequestContext) {
        fmt.Println("[Logger] 记录请求日志")
        c.Next(ctx)
    }
}
