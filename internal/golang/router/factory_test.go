package router

import (
    "log"
    "testing"

    "github.com/gin-gonic/gin"
)

func TestFactory(t *testing.T) {
    e := gin.New()

    rf := NewFactory(e)

    userRouterFactory := rf.AddPath("/api").
        AppendBaseMiddleware(Auth()).
        AppendInnerMiddleware(gin.Logger(), gin.Recovery()).
        AppendOpenMiddleware(Logger())

    userRouters := userRouterFactory.Clone().
        AddPath("/user").
        Routers()

    // path: /api/user/inner/list
    userRouters.InnerRouter.GET("/list", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Inner User List",
        })
    })

    // curl -X GET http://localhost:8080/api/user/inner/list
    log.Fatal(e.Run(":8080"))
}
