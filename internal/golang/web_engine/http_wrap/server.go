package main

import "net/http"

type Server interface {
	Routable

	Start(address string) error
}

type sdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

func (s *sdkHttpServer) Route(method string, pattern string, handleFunc handleFunc) {
	s.handler.Route(method, pattern, handleFunc)
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/",
		func(writer http.ResponseWriter, request *http.Request) {
			s.root(NewContext(writer, request))
		})
	return http.ListenAndServe(address, nil)
}

func NewHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBasedOnTree()

	// init filter
	var root Filter = handler.ServeHTTP
	for i := len(builders); i >= 0; i-- {
		builder := builders[i]
		root = builder(root)
	}

	return &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}

func SignUp(ctx *Context) {
	// doing sign up
}
