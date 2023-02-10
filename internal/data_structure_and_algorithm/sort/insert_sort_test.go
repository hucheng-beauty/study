package sort

import "testing"

func TestInsertSort(t *testing.T) {
	arr := []int{5, 6, 1, 2, 3, 4}
	t.Log(arr)
	InsertSort(arr)
	t.Log(arr)
}
