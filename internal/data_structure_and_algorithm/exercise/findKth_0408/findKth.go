package main

func quick(nums []int, left, right int) int {
    if left == right {
        return nums[left]
    }

    index := nums[right]
    i := left

    for j := left; j < right; j++ {
        if nums[j] > index {
            nums[i] = nums[j]
            nums[j] = nums[i]
            i++
        }
    }
    nums[i] = nums[right]
    nums[right] = nums[i]
    return i
}

func quickSort(nums []int, left, right, k int) int {
    if left == right {
        return nums[left]
    }

    index := quick(nums, left, right)

    if k == index {
        return nums[k]
    } else if k < index {
        return quickSort(nums, left, index-1, k)
    } else {
        return quickSort(nums, index+1, right, k)
    }
}

func findKth(nums []int, length, k int) int {
    return quickSort(nums, 0, length-1, k)
}
