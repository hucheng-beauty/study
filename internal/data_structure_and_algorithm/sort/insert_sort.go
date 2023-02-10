package sort

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
