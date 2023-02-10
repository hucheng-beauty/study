package sort

// BubbleSort 时间复杂度:最好 O(n)、最坏 O(n^2)、平均 O(n^2); 空间复杂度:O(1); 原地性:是; 稳定性: 是
// 思路: 比较当前值与下一个值,大的往后移动,一趟冒泡后,至少有一个是拍好序的
// 优化: 若一趟排序之后无交换数据,则证明都是有序的;则可以通过增加一个 flag 来标记是否有数据交换
func BubbleSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	for i := 0; i < len(arr); i++ {
		flag := false
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				flag = true
			}
		}
		if !flag {
			break
		}
	}
}
