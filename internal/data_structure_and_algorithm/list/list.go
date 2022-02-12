package list

type listNode struct {
	Data int
	Next *listNode
}

func (ln *listNode) Traverse() {
	markListNode := ln
	for markListNode != nil {
		markListNode = markListNode.Next
	}
}

func (ln *listNode) FindValue(value int) *listNode {
	mln := ln
	for mln != nil {
		if mln.Data == value {
			return mln
		}
		mln = mln.Next
	}
	return nil
}

func (ln *listNode) HeadInsert(value int) {
	newNode := &listNode{value, nil}
	if ln == nil {
		ln = newNode
	} else {
		newNode.Next = ln.Next
		ln.Next = newNode
	}
}

func (ln *listNode) TailInsert(value int) {
	newNode := &listNode{value, nil}
	if ln == nil {
		ln.Next = newNode
	} else {
		mln := ln
		// find tail node
		for mln != nil {
			mln = mln.Next
		}
		mln = newNode
	}
}

// 尾部结点
var tail = &listNode{}

func (ln *listNode) GreatTailInsert(value int) {
	newNode := &listNode{value, nil}
	if ln == nil {
		ln.Next = newNode
	} else {
		tail.Next = newNode
	}
}

func (ln *listNode) LocationInsert(p *listNode, value int) {
	newNode := &listNode{value, nil}
	if p == nil || ln == nil {
		return
	} else {
		mln := p

		newNode.Next = mln.Next
		p.Next = newNode
	}
}

/*
	1.找到前驱结点
	2.删除结点
*/
// Delete 删除给定结点
func (ln *listNode) Delete(p *listNode) {
	if p == nil || ln == nil {
		return
	}
	var prev *listNode
	mln := ln
	for mln != nil {
		if mln == p {
			break
		}
		// 记录前驱结点
		prev = mln
		mln = mln.Next
	}
	if mln == nil {
		return
	}
	if prev == nil {
		ln = ln.Next
	} else {
		prev.Next = prev.Next.Next
	}

}

// DeleteNextNode 删除给定结点之后的结点
func (ln listNode) DeleteNextNode(p *listNode) {
	if p == nil || p.Next == nil {
		return
	}
	p.Next = p.Next.Next
}
