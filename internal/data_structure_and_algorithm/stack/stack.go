package stack

type ArrayStack struct {
	items []int // 元素
	count int   // 栈中元素的个数
	n     int   // 栈的大小
}

func NewArrayStack(n int) *ArrayStack {
	return &ArrayStack{
		items: make([]int, n),
		count: 0,
		n:     n,
	}
}

// Pop 弹出栈中元素
func (as *ArrayStack) Pop() int {
	if as.count <= 0 {
		return -1
	}

	out := as.items[as.count-1]
	as.count--
	return out
}

// Push 向栈中推入元素
func (as *ArrayStack) Push(item int) bool {
	if as.count == as.n {
		return false
	}

	as.items[as.count] = item
	as.count++
	return true
}

// Peek 获取栈中的栈顶元素
func (as *ArrayStack) Peek() int {
	if as.count == 0 {
		return -1
	}
	return as.items[as.count-1]
}

type ListNode struct {
	Value int
	Next  *ListNode
}

type ListStack struct {
	head *ListNode
}

func (ls *ListStack) Push(value int) bool {
	newNode := new(ListNode)
	newNode.Value = value

	ls.head.Next = newNode
	return true
}

// TODO: return -1 有 bug, Pop、Peek

func (ls *ListStack) Pop() int {
	if ls.head == nil {
		return -1
	}

	out := ls.head.Value
	ls.head = ls.head.Next
	return out
}

func (ls *ListStack) Peek() int {
	if ls.head == nil {
		return -1
	}
	return ls.head.Value
}
