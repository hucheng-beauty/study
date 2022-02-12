package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewHttpServer("test-server")
	server.Route(http.MethodPost, "/user/signup", SignUp)
	log.Fatal(server.Start(":8080"))
}
