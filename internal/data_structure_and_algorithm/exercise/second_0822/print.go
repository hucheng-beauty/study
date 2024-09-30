package second_0822

import "fmt"

/*
   func print(n int) // 按照下面格式和规律打印1~n行
   1
   1 1
   1 2 1
   1 3 3 1
   1 4 6 4 1
   1 5 10 10 5 1
   1 6 15 20 15 6 1
   …
*/

func Print(n int) {
	for i := 0; i < n; i++ {
		num := 1
		for j := 0; j <= i; j++ {
			fmt.Print(num, " ")
			num = num * (i - j) / (j + 1)
		}
		fmt.Println()
	}
}
