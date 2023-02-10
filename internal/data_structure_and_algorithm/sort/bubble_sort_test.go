package sort

import "testing"

func TestBubbleSort(t *testing.T) {
	arr := []int{5, 6, 1, 2, 3, 4}
	t.Log(arr)
	BubbleSort(arr)
	t.Log(arr)
}
