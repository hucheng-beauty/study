package array

import "testing"

func TestMainer(t *testing.T) {
	arr := [3]int{1, 2, 3}
	t.Logf("arr[0]:%d\n", arr[0])

	// traverse 1
	for i := 0; i < len(arr); i++ {
		t.Log(arr[i])
	}

	// traverse 2
	for index, value := range arr {
		t.Logf("index:%d, value:%d\n", index, value)
	}
}
