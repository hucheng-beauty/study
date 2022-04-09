package main

import "fmt"

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

// InsertSort 时间复杂度:最好 O(n)、最坏 O(n^2)、平均 O(n^2); 空间复杂度:O(1); 原地性:是; 稳定性: 是
// 思路: 分为有序和无序;类似打牌的思路,没接到一张牌一次比较并插入合适的位置
func InsertSort(arr []int) {
	if len(arr) <= 1 {
		return
	}

	for i := 1; i < len(arr); i++ {
		mv := arr[i]
		j := i - 1

		for ; j >= 0; j-- {
			if arr[j] > mv { // 比较
				arr[j+1] = arr[j] // 移动
			} else {
				break
			}
		}
		arr[j+1] = mv // 插入数据
	}
}

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

// MergeSort 时间复杂度: O(NlogN)、平均 O(n); 空间复杂度:O(n); 原地性:; 稳定性:
// 递推公式: MergeSort(MergeSort(start, middle), MergeSort(middle + 1, end)); middle = (start + end) / 2
// 终止条件: start >= end 不用在继续分解;只有一个元素或者没有元素结束
// 思路: 将数据分为前后俩部分,然后分别排序;最后将排好序的俩部分合并在一起
func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	middle := len(arr) / 2
	left := MergeSort(arr[:middle])
	right := MergeSort(arr[middle:])

	out := make([]int, 0)
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		// 双指针比较
		if left[i] > right[j] {
			out = append(out, right[j])
			j++
		} else {
			out = append(out, left[i])
			i++
		}
	}

	// 若其中一个数组没有值,则将另一个下表以后的所有数据添加至 out
	if i < len(left) {
		out = append(out, left[i:]...)
	}
	if j < len(right) {
		out = append(out, right[j:]...)
	}

	return out
}

// QuickSort 快速排序 时间复杂度: 最好 O(NlogN)、最坏 O(n^2)、平均 O(NlogN); 空间复杂度:最好:O(logN)、最坏:O(n)、平均:O(logN); 原地性:是; 稳定性:否
// 递推公式: QuickSort(start, end) = partition(start, end) + QuickSort(start, middle) + QuickSort(middle + 1, end); middle = (start + end) / 2
// 终止条件: start >= end 不用在继续分解;只有一个元素或者没有元素结束\
// 思路: 分为三部分,一个数字,前面是比这个数小,后面的比这个数大
func QuickSort(arr []int, left, right int) {
	if left >= right {
		return
	}

	// 选取一个比较值
	markData := arr[left]

	// [left, j]: 小于 markData
	// [j, i]: 大于 markData
	// [i, right]: 未处理的数据
	j := left
	for i := left; i < right; i++ {
		if arr[i] < markData {
			j++
			arr[j], arr[i] = arr[i], arr[j]
		}
	}
	arr[left], arr[j] = arr[j], arr[left]
	QuickSort(arr, left, j)
	QuickSort(arr, j+1, right)
}

func main() {
	//arr := []int{4, 5, 1, 2, 6, 3}
	arr := []int{5, 6, 1, 2, 3, 4}
	fmt.Println(arr)
	QuickSort(arr, 0, len(arr))
	fmt.Println(arr)
}
