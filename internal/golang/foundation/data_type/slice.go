package main

import (
	"fmt"
)

/*
	slice:
		对数组的一种抽象,可以自动扩容(小于1024,2 倍扩容;大于 1024,1.5倍扩容,同时进行内存对齐)
	注意事项:
		俩个切片会相互影响
		当 append 元素时,发生扩容时,对原数组不产生影响

*/

/*
	// src/runtime/slice.go
	type slice struct {
		array unsafe.Pointer // 指向数组的指针
		len   int            // 长度
		cap   int            // 容量
	}
*/

func sliceTest() {
	slice := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s1 := slice[2:5] // 2, 3, 4
	s2 := s1[2:6:7]  // 4, 5, 6, 7

	s2 = append(s2, 100) // 4, 5, 6, 7, 100
	s2 = append(s2, 200) // 4, 5, 6, 7, 100, 200

	s1[2] = 20 // 2, 3, 20

	fmt.Println(s1)               // 2, 3, 20
	fmt.Println(len(s1), cap(s1)) // 3, 8

	fmt.Println(s2)
	fmt.Println(len(s2), cap(s2))

	fmt.Println(slice) // 0, 1, 2, 3, 20, 5, 6, 7, 100, 9
}
