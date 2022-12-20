package queue

/*
	俩种标记 循环队列满和空的方法
		1.用 count 记录存储在队列中的数据个数
			count == 0 空, count == n 满

		2.不使用 count
			head == tail ==> 空, (tail + 1) % n == head ==> 满
*/

// ArrayQueue circular queue
type ArrayQueue struct {
	head  int
	tail  int
	n     int
	items []string
}

func NewArrayQueue(n int) *ArrayQueue {
	return &ArrayQueue{
		head:  0,
		tail:  0,
		n:     n,
		items: make([]string, n),
	}
}

func (aq *ArrayQueue) Enqueue(item string) bool {
	if (aq.tail+1)%aq.n == aq.head {
		return false
	}
	aq.items[aq.tail] = item
	aq.tail = (aq.tail + 1) % aq.n
	return true
}

func (aq *ArrayQueue) Dequeue() string {
	if aq.head == aq.tail {
		return ""
	}
	aq.head = (aq.head + 1) % aq.n
	return aq.items[aq.tail]
}

type ListNode struct {
	Next  *ListNode
	Value string
}

type ListQueue struct {
	head *ListNode
	tail *ListNode
}

func NewListQueue() *ListQueue {
	return &ListQueue{
		head: nil,
		tail: nil,
	}
}

func (lq *ListQueue) Enqueue(value string) {
	newNode := &ListNode{Value: value}

	// 处理空链表
	if lq.tail == nil {
		lq.head = newNode
		lq.tail = newNode
	} else {
		// 更新 tail 指针
		lq.tail.Next = newNode
		lq.tail = newNode
	}
}

func (lq *ListQueue) Dequeue() string {
	if lq.head == nil {
		return ""
	}

	rt := lq.head.Value
	lq.head = lq.head.Next

	// 处理队列只有一个值
	if lq.head == nil {
		lq.tail = nil
	}
	return rt

}
