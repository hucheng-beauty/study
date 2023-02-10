package main

import (
    "testing"

    "study/pkg/crawler/engine"
    "study/pkg/crawler_distributed/rpcsupport"
)

func TestItemSaverService(t *testing.T) {
    const host = ":1234"
    const index = "dating_profile"
    // start ItemSaverServer
    go serveRpc(host, index)

    // start ItemSaverClient
    client, err := rpcsupport.NewClient(host)
    if err != nil {
        panic(err)
    }

    // call Save
    item := engine.Item{}

    result := ""
    err = client.Call("ItemSaverService.Save", item, &result)

    if err != nil || result != "ok" {
        t.Errorf("result: %s, error:%s", result, err)
    }

}
