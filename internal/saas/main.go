package main

import (
    "study/internal/saas/initialize"
    "study/internal/saas/service/trigger/router"

    "github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
    // trigger service

    err := initialize.Config()
    if err != nil {
        panic(err)
    }

    err = initialize.MySQL()
    if err != nil {
        panic(err)
    }

    // initialize others

    engine := server.New()
    router.Handler(engine)
    err = engine.Run()
    if err != nil {
        return
    }
}
