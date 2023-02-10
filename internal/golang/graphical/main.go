package main

import (
	"log"
	"net/http"

	"github.com/arl/statsviz"
)

// 图形化的 runtime

func main() {
	// Graphical Go Runtime
	_ = statsviz.RegisterDefault()
	log.Println(http.ListenAndServe(":6060", nil))
}
