package sort

import (
	"testing"
)

func TestMainer(t *testing.T) {
	//arr := []int{4, 5, 1, 2, 6, 3}
	arr := []int{5, 6, 1, 2, 3, 4}
	t.Log(arr)
	QuickSort(arr, 0, len(arr))
	t.Log(arr)
}
