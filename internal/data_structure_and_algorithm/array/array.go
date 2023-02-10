package array

import "fmt"

/*
	数组
		固定大小、连续的一段内存,不可自动扩容
		按照下标访问的时间复杂度为 O(1)
		下标从零开始
*/

func Array() {
	arr := [3]int{1, 2, 3}

	// traverse 1
	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}

	// traverse 2
	for index, value := range arr {
		fmt.Printf("index:%d, value:%d\n", index, value)
	}
}
