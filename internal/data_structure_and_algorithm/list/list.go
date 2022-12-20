package list

/*
	链表
		头节点: 指的是链表第一个节点
		头指针: 指的是指向头节点的指针
*/

type ListNode struct {
	Data int
	Next *ListNode
}

// Traverse 遍历
func (ln *ListNode) Traverse() {
	for mln := ln; mln != nil; mln = mln.Next {
		println(mln.Data)
	}
}

// FindValue 查询链表中的值
func (ln *ListNode) FindValue(value int) *ListNode {
	for mln := ln; mln != nil; mln = mln.Next {
		if mln.Data == value {
			return mln
		}
	}
	return nil
}

// HeadInsert 头插
func (ln *ListNode) HeadInsert(value int) *ListNode {
	nn := new(ListNode)
	nn.Data = value
	nn.Next = ln
	ln = nn
	return ln
}

// TailInsert 尾插
func (ln *ListNode) TailInsert(value int) {
	nn := new(ListNode)
	nn.Data = value
	if ln == nil {
		ln.Next = nn
	}

	// find tail node and tail insert
	mln := ln
	for ; mln.Next != nil; mln = mln.Next {
	}
	mln.Next = nn
}

// 尾部结点
var tail *ListNode = nil

// WithTailNodeOfTailInsert 带尾部节点的尾插
func (ln *ListNode) WithTailNodeOfTailInsert(value int) {
	newNode := &ListNode{value, nil}
	if ln == nil {
		ln.Next = newNode
	} else {
		tail.Next = newNode
	}
}

func (ln *ListNode) LocationInsert(p *ListNode, value int) {
	if p == nil || ln == nil {
		return
	} else {
		mln := p

		newNode := &ListNode{value, nil}
		newNode.Next = mln.Next
		p.Next = newNode
	}
}

// Delete 删除给定结点
func (ln *ListNode) Delete(p *ListNode) {
	if p == nil || ln == nil {
		return
	}

	var pre *ListNode
	mln := ln
	for ; mln != nil; mln = mln.Next {
		if mln == p { // 找到要删除的节点
			break
		}
		pre = mln // 记录要删除节点的前驱结点
	}

	if mln == nil {
		return
	}
	if pre == nil {
		println("hello world!")
		ln = ln.Next
	} else {
		pre.Next = pre.Next.Next
	}
}

// DeleteNextNode 删除给定结点之后的结点
func (ln *ListNode) DeleteNextNode(p *ListNode) {
	if p == nil || p.Next == nil {
		return
	}
	p.Next = p.Next.Next
}

// Reverse 反转链表
func (ln *ListNode) Reverse() *ListNode {
	if ln == nil {
		return nil
	}

	// 前驱
	var pre *ListNode = nil
	cur := ln // 指向头节点
	for cur != nil {
		tmp := cur.Next // 防止链表断开找不到
		cur.Next = pre  // 改变指针指向
		pre = cur       // 向后移动 pre
		cur = tmp       // 向后移动 cur
	}
	return pre
}
