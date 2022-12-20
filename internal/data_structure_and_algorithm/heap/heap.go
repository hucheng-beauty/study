package heap

import "fmt"

/*
	堆
		1.必须为完全二叉树
		2.堆中的每个节点必须大于等于(或小于等于)其他子树中每个节点的值
	大顶堆 & 小顶堆
		每个节点大于等于子树中每个节点的值,则称为大顶堆
		每个节点小于等于子树中每个节点的值,则称为小顶堆
	存储
		适合用数组
		一个节点的下标为 i, 他的子节点的下标为 i*2 和 i*2+1
		root 节点的下标从 1 开始

*/

// Heap 大顶堆
type Heap struct {
	a   []int
	len int // 堆中实现存在的长度
	cap int // 堆的容量
}

func NewHeap(cap int) *Heap {
	return &Heap{
		a:   make([]int, cap+1),
		len: 0,
		cap: cap,
	}
}

// Insert 往堆中添加元素
func (h *Heap) Insert(data int) {
	if h.len >= h.cap { // 堆满了
		fmt.Println("堆满了")
		return
	}

	// 尾部追加数据
	h.len++
	h.a[h.len] = data

	i := h.len
	for i/2 > 0 && h.a[i] > h.a[i/2] { // 自下向上堆化
		h.a[i], h.a[i/2] = h.a[i/2], h.a[i]
		i = i / 2
	}
}

// Top 获取堆顶元素
func (h *Heap) Top() int {
	if h.len == 0 {
		fmt.Println("堆空了")
		return 0
	}
	return h.a[1]
}

// Update 删除具体位置上的元素
func (h *Heap) Update(index, value int) {
	// TODO: implement the method of Update
	if 0 < index && index <= h.len {
		if h.a[index] == value {
			return
		}

		// 自下向上堆化
		if h.a[index] > value {
			for {
				if h.a[index] > h.a[(index-1)/2] {

				}
			}

		} else { // 自上向下堆化

		}
	}
	return
}

// RemoveTop 删除堆顶元素
func (h *Heap) RemoveTop() {
	if h.len == 0 { // 堆空了
		return
	}

	// 移除堆顶值
	h.a[1] = h.a[h.len]
	h.len--

	// 自上而下堆化
	i := 1
	for i*2 <= h.len && i*2+1 <= h.len {
		maxPos := i
		if h.a[i] < h.a[i*2] {
			maxPos = i * 2
		}
		if h.a[maxPos] < h.a[i*2+1] {
			maxPos = i*2 + 1
		}
		if maxPos == i {
			break
		}
		h.a[i], h.a[maxPos] = h.a[maxPos], h.a[i]
		i = maxPos
	}
}
