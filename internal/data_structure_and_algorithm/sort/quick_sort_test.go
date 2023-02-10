package sort

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	arr := []int{5, 6, 1, 2, 3, 4}
	t.Log(arr)
	QuickSort(arr, 0, len(arr))
	t.Log(arr)
}
