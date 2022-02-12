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

func (aq *ArrayQueue) enqueue(item string) bool {
	if (aq.tail+1)%aq.n == aq.head {
		return false
	}
	aq.items[aq.tail] = item
	aq.tail = (aq.tail + 1) % aq.n
	return true
}

func (aq *ArrayQueue) dequeue() string {
	if aq.head == aq.tail {
		return ""
	}
	return aq.items[aq.tail]
}
