package sort

import "math/rand"

// CountingSort 计数排序
func CountingSort(a []int, n int) []int {
	if len(a) <= 1 {
		return a
	}
	// 确定数据的范围
	max := a[0]
	for i := 1; i < n; i++ {
		if max < a[i] {
			max = a[i]
		}
	}

	// 确定每个元素的个数
	c := make([]int, max+1) // 下标为(0, max)
	for i := 0; i < n; i++ {
		c[a[i]]++
	}

	// 依次累加
	for i := 1; i <= max; i++ {
		c[i] = c[i-1] + c[i]
	}

	/*
		a = [1, 4, 5, 3, 2, 0, 1, 3, 4] // 初始数组
			 0  1  2  3  4  5  6  7  8

		c = [1, 2, 1, 2, 2, 1] // 每个元素出现的次数
			 0  1  2  3  4  5

		c = [0, 2, 3, 4, 6, 8] // 依次累加
			 0  1  2  3  4  5

		r = [0, 1, 1, 2, 3, 3, 4, 4, 5] // 核心步骤
			 0  1  2  3  4  5  6  7  8
	*/

	r := make([]int, n) // 临时数组,存储排序之后的结果
	for i := n - 1; i >= 0; i-- {
		// 核心步骤
		index := c[a[i]] - 1
		r[index] = a[i]
		c[a[i]]--
	}

	// 将结果拷贝给 a 数组
	for i := 0; i < n; i++ {
		a[i] = r[i]
	}
	return a
}

// RandInt returns a non-negative pseudo-random int from the default Source.
func RandInt() int {
	return rand.Int()
}
