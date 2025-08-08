package router

import (
    "context"
    "log"
    "testing"

    "github.com/cloudwego/hertz/pkg/app"
    "github.com/cloudwego/hertz/pkg/app/server"
)

func TestFactory(t *testing.T) {
    e := server.New()

    rf := NewFactory(e)

    userRouterFactory := rf.AddPath("/api").
        AppendBaseMiddleware(Auth()).
        AppendInnerMiddleware(Logger()).
        AppendOpenMiddleware()

    userRouters := userRouterFactory.Clone().
        AddPath("/user").
        Groups()

    // path: /api/user/inner/list
    userRouters.InnerRouter.GET("/list", func(ctx context.Context, c *app.RequestContext) {
        c.JSON(200, "hello world")
    })

    // curl -X GET http://localhost:8080/api/user/inner/list
    log.Fatal(e.Run())
}
