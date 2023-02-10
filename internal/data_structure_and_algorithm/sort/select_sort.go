package sort

// SelectSort 时间复杂度:最好 O(n)、最坏 O(n^2)、平均 O(n^2); 空间复杂度:O(1); 原地性:否; 稳定性: 是
// 思路: 分为有序和无序;与插入排序不一样的是,选择排序是在无序数据中找最小的一个,然后交换数据
func SelectSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	for i := 0; i < len(arr)-1; i++ {
		minIndex := i
		for j := i; j < len(arr); j++ {
			if arr[j] < arr[minIndex] { // 找 minValue
				minIndex = j
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i] // 交换元素
	}
}
