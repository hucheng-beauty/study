package third_0913

/*
   数组 [1,2,3,2,1]
*/

func findCountOne(arr []int) int {
	m := map[int]int{}

	// count
	for _, v1 := range arr {
		m[v1]++
	}

	// find
	for _, v2 := range arr {
		if m[v2] == 1 {
			return v2
		}
	}
	return -1
}
