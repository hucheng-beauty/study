package router

import (
    "study/internal/saas/service/desicion/router"
    "study/internal/saas/service/execution/action"

    "github.com/cloudwego/hertz/pkg/app/server"
)

func Handler(engine *server.Hertz) {
    rf := router.NewFactory(engine)
    saasRouterFactory := rf.AddPath("/saas").
        AppendBaseMiddleware().
        AppendInnerMiddleware().
        AppendOpenMiddleware()

    userRouters := saasRouterFactory.Clone().
        AddPath("/user").
        Groups()
    userRouters.InnerRouter.POST("/CreateApiKey", action.SelectAction(
        action.NewVersionOne("CreateApiKey"),
    ))
}
