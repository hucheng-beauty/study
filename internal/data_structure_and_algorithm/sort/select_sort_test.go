package sort

import "testing"

func TestSelectSort(t *testing.T) {
	arr := []int{5, 6, 1, 2, 3, 4}
	t.Log(arr)
	SelectSort(arr)
	t.Log(arr)

}
