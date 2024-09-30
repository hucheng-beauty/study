package second_0822

// 快速排序

// 1 3 2 5 4
// l 1
// index := 1
// r 4

func QuickSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}

	index := arr[0]
	var l, g []int
	for _, v := range arr[1:] {
		if v <= index {
			l = append(l, v)
		} else {
			g = append(g, v)
		}
	}
	l = QuickSort(l)
	g = QuickSort(g)

	return append(append(l, index), g...)
}
