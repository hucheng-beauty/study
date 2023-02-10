package sort

// QuickSort 快速排序 时间复杂度: 最好 O(NlogN)、最坏 O(n^2)、平均 O(NlogN);
// 空间复杂度:最好:O(logN)、最坏:O(n)、平均:O(logN); 原地性:是; 稳定性:否
// 递推公式: QuickSort(start, end) = partition(start, end) + QuickSort(start, middle) + QuickSort(middle + 1, end);
// middle = (start + end) / 2
// 终止条件: start >= end 不用在继续分解;只有一个元素或者没有元素结束\
// 思路: 分为三部分,一个数字,前面是比这个数小,后面的比这个数大
func QuickSort(arr []int, low, high int) {
	if low < high {
		// 将数组划分为两个子数组，并获取划分点的索引
		pivotIndex := partition(arr, low, high)

		// 递归地对划分的子数组进行快速排序
		QuickSort(arr, low, pivotIndex-1)
		QuickSort(arr, pivotIndex+1, high)
	}
}

// partition 划分函数,用于将数组划分为两个子数组
func partition(arr []int, low, high int) int {
	// 选择最后一个元素作为划分点
	pivot := arr[high]

	// 定义左指针和右指针
	i := low - 1

	// 遍历数组，将小于划分点的元素放在左侧，大于划分点的元素放在右侧
	for j := low; j <= high-1; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// 将划分点放在正确的位置上
	arr[i+1], arr[high] = arr[high], arr[i+1]

	return i + 1
}

func QuickSorting(data []int, start, end int) {
	if start < end {
		base := data[start] // 获取基准值
		left := start
		right := end

		for left < right {
			// 从后往前寻找比基准值小的数,若大于基准值则进行索引前移
			for left < right && data[right] >= base {
				right--
			}
			if left < right {
				data[left] = data[right]
				left++
			}

			// 从前往后寻找比基准值大的数,若小于基准值则进行索引前移
			for left < right && data[left] <= base {
				left++
			}
			if left < right {
				data[right] = data[left]
				right--
			}
		}

		data[left] = base

		QuickSorting(data, start, left-1)
		QuickSorting(data, left+1, end)
	}
}
