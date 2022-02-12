package main

import "strings"

type node struct {
	path     string
	children []*node
	handler  handleFunc
}

func newNode(path string) *node {
	return &node{
		path:     path,
		children: make([]*node, 0, 2),
	}
}

// HandlerBasedOnTree 路由树
type HandlerBasedOnTree struct {
	root *node
}

func (h *HandlerBasedOnTree) Route(method string, pattern string, handleFunc handleFunc) {
	// 预处理 & 切割
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	// 寻找路由
	cur := h.root
	for index, path := range paths {
		matchNode, found := h.findMatchChild(cur, path)
		if found { // 递归查找
			cur = matchNode
		} else { // 没找到创建 node
			h.createSubNode(cur, paths[index:], handleFunc)
			return
		}
	}
}

func (h *HandlerBasedOnTree) findMatchChild(root *node, path string) (*node, bool) {
	for _, child := range root.children {
		if child.path == path {
			return child, true
		}
	}
	return nil, false
}

func (h *HandlerBasedOnTree) createSubNode(root *node, paths []string, handleFn handleFunc) {
	cur := root
	for _, path := range paths {
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}

	cur.handler = handleFn
}

func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	panic("implement me")
}

func (h *HandlerBasedOnTree) findRouter() {

}

func NewHandlerBasedOnTree() Handler {
	return &HandlerBasedOnTree{
		root: nil,
	}
}
