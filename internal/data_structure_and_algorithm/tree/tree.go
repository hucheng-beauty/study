package tree

import (
	"fmt"
)

type Node struct {
	Data  int
	Left  *Node
	Right *Node
}

func (root *Node) Insert(data int) {
	if root == nil {
		root = &Node{Data: data}
		return
	}

	mr := root
	for mr != nil {
		if mr.Data > data {
			if mr.Left == nil {
				nn := &Node{Data: data}
				mr.Left = nn
				return
			}
			mr = mr.Left
		} else {
			if mr.Right == nil {
				nn := &Node{Data: data}
				mr.Right = nn
				return
			}
			mr = mr.Right
		}
	}
}

func insert(root *Node, data int) {
	//preHandle(root, data)

	if root == nil {
		root = &Node{Data: data}
		return
	}

	if root.Data > data {
		if root.Left == nil {
			nn := &Node{Data: data}
			root.Left = nn
		} else {
			insert(root.Left, data)
		}
	} else { // root.Data <= data
		if root.Right == nil {
			nn := &Node{Data: data}
			root.Right = nn
		} else {
			insert(root.Right, data)
		}
	}
}

func (root *Node) InsertRecursion(data int) {
	insert(root, data)
}

func (root *Node) Find(data int) *Node {
	mr := root
	for mr != nil {
		if mr.Data == data {
			fmt.Printf("find value: %d\n", data)
			return mr
		} else if mr.Data > data {
			mr = mr.Left
		} else {
			mr = mr.Right
		}
	}
	fmt.Printf("can not find value: %d\n", data)
	return nil
}

func find(root *Node, data int) *Node {
	if root == nil {
		return nil
	}

	if root.Data == data {
		return root
	} else if root.Data > data {
		return find(root.Left, data)
	} else {
		return find(root.Right, data)
	}
}

func (root *Node) FindRecursion(data int) *Node {
	return find(root, data)
}

func (root *Node) Delete(data int) {
	var p *Node = nil  // 标记要删除的节点
	var pp *Node = nil // 标记要删除节点的父节点

	p = root
	for p != nil && p.Data != data {
		pp = p
		if p.Data > data {
			p = p.Left
		} else {
			p = p.Right
		}
	}

	// 没有找到
	if p == nil {
		return
	}

	// 要删除的节点有俩个节点
	var minP *Node = nil
	var minPP *Node = nil // minPP 表示 minP 的父节点

	if p.Right != nil && p.Left != nil {
		minP = p.Right
		minPP = p

		// 查找右子树中最小的节点
		for minP.Left != nil {
			minPP = minP
			minP = minP.Left
		}

		// 将最小节点的数据移到 p 和 pp
		p.Data = minP.Data
		p = minP
		pp = minPP
	}

	// 要删除的节点为叶子节点或仅有一个子节点
	var child *Node = nil // 查找待删除节点 p 的子节点
	if p.Right != nil {
		child = p.Right
	} else if p.Left != nil {
		child = p.Left
	}

	// 删除节点
	if pp == nil { // 删除的是根节点
		root = child
	} else if pp.Left == p {
		pp.Left = child
	} else {
		pp.Right = child
	}
}

func preOrder(root *Node) {
	if root == nil {
		return
	}
	fmt.Println(root.Data)
	preOrder(root.Left)
	preOrder(root.Right)
}

func (root *Node) PreOrder() {
	preOrder(root)
}

func inOrder(root *Node) {
	if root == nil {
		return
	}
	inOrder(root.Left)
	fmt.Println(root.Data)
	inOrder(root.Right)
}

func (root *Node) InOrder() {
	inOrder(root)
}

func postOrder(root *Node) {
	if root == nil {
		return
	}
	postOrder(root.Left)
	postOrder(root.Right)
	fmt.Println(root.Data)
}

func (root *Node) PostOrder() {
	postOrder(root)
}
