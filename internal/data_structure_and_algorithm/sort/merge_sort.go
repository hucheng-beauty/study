package sort

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
