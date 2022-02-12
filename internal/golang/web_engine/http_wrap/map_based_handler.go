package main

import "net/http"

type handleFunc func(c *Context)

type Routable interface {
	Route(method string, pattern string, handleFunc handleFunc)
}

type Handler interface {
	Routable

	ServeHTTP(c *Context)
}

type HandlerBasedMap struct {
	handles map[string]handleFunc
}

// Route 绑定路由和 handleFunc
func (h *HandlerBasedMap) Route(method string, pattern string,
	handleFunc handleFunc) {
	h.handles[h.key(method, pattern)] = handleFunc
}

// ServeHTTP 处理具体的路由的handleFunc
func (h *HandlerBasedMap) ServeHTTP(c *Context) {
	key := h.key(c.R.Method, c.R.URL.Path)

	// handler 校验
	if handler, ok := h.handles[key]; ok {
		handler(c)
	} else {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("Not Found handler"))
	}
}

func (h *HandlerBasedMap) key(method string, pattern string) string {
	return method + "#" + pattern
}

// 	确保 HandlerBasedMap 一定实现了 Handler 接口
var _ Handler = &HandlerBasedMap{}

func NewHandleBasedMap() Handler {
	return &HandlerBasedMap{
		handles: make(map[string]handleFunc),
	}
}
