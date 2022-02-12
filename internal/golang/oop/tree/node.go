package main

/*
	golang 中所有的参数都是值传递
	值接收者为 golang 特有
*/

/*
	包:
		如何扩充系统类型或者别人的类型
			1.定义别名:最简单
			2.使用组合:最常用
			3.使用内嵌:需要省很多代码时候使用
*/
type treeNode struct {
	value int
	left  *treeNode
	right *treeNode
}

func CreateTreeNode(value int) *treeNode {
	return &treeNode{value: value}
}

func main() {
	var root treeNode

	root = treeNode{value: 3}
	root.left = &treeNode{}
	root.right = &treeNode{value: 5}
}
