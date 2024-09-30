package queue

/*
	连连消:删除连续重读的字符
		举例: "abbbc"   ==> "ac"
		     "abbbabc" ==> "c"
	解法:
		使用栈解决,栈里面存储不能被消除的元素
*/

// Item 连连消中的元素,和改元素出现的次数
type Item struct {
	Name  string
	Count int
}

func Remove() {
	// 1.创建一个栈对象
	// stack := make([]Item, 0)
	// 遍历给定字符串,入栈的同时统计出现的次数
}

type Stack struct {
	Value []Item
}

func (s *Stack) Enqueue() {

}

func (s *Stack) Dequeue() {

}
