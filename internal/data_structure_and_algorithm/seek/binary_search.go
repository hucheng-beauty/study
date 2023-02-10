package seek

func BinarySearch(arr []int, target int) int {
    low := 0
    high := len(arr) - 1

    for low <= high {
        mid := (high-low)/2 + low
        if arr[mid] == target {
            return mid // 找到返回下标
        } else if arr[mid] > target {
            high = mid - 1
        } else {
            low = mid + 1
        }
    }
    return -1 // 没找到返回-1
}

func BinarySearchRecursion(arr []int, low, high int, target int) int {
    if low > high {
        return -1
    }

    mid := (high-low)/2 + low // 避免溢出
    if arr[mid] == target {
        return mid
    } else if arr[mid] > target {
        return BinarySearchRecursion(arr, low, mid-1, target)
    } else {
        return BinarySearchRecursion(arr, mid+1, high, target)
    }
}
