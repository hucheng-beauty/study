package sort

import (
	"testing"
)

func TestMergeSort(t *testing.T) {
	arr := []int{5, 6, 1, 2, 3, 4}
	t.Log(arr)
	MergeSort(arr)
	t.Log(arr)
}
