package main

import (
	"fmt"
	"strings"
	"time"
)

// StringJoin join 字符串拼接
func StringJoin() {
	t1 := time.Now().UnixNano() / 1e6
	fmt.Println("t1 ==> ", t1)

	s := make([]string, 2)
	s = append(s, "hello")
	for i := 0; i < 500000; i++ {
		s = append(s, "world")
	}
	strings.Join(s, "")

	t2 := time.Now().UnixNano() / 1e6
	fmt.Println("t2 ==> ", t2)
	fmt.Println("t2 - t1 ==> ", t2-t1)
}
