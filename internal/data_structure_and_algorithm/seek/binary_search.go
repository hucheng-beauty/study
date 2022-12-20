package seek

// BinarySearch 二分查找，找到返回下标，没找到返回-1
func BinarySearch(arr []int, target int) int {
	low := 0
	high := len(arr) - 1

	for low <= high {
		//mid := (low + high) / 2
		mid := (high-low)/2 + low
		if arr[mid] == target {
			return mid
		} else if arr[mid] > target {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}
	return -1
}

// BinarySearchRecursion 递归实现二分查找
func BinarySearchRecursion(arr []int, low, high int, target int) int {
	if low > high {
		return -1
	}

	mid := (low + high) / 2
	if arr[mid] == target {
		return mid
	} else if arr[mid] > target {
		return BinarySearchRecursion(arr, low, mid-1, target)
	} else {
		return BinarySearchRecursion(arr, mid+1, high, target)
	}
}
