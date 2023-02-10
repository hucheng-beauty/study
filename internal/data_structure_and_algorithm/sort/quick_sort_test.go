package sort

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	arr := []int{6, 3, 1, 2, 5, 4}
	t.Log(arr)
	QuickSorting(arr, 0, len(arr)-1)
	t.Log(arr)
}
