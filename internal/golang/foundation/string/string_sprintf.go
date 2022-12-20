package string

import (
	"fmt"
	"time"
)

// StringSprintf SPrintf字符串拼接
func StringSprintf() {
	t1 := time.Now().UnixNano() / 1e6
	fmt.Println("t1 ==> ", t1)

	s := "hello"
	for i := 0; i < 500000; i++ {
		s = fmt.Sprintf("%s%s", s, "world")
	}

	t2 := time.Now().UnixNano() / 1e6
	fmt.Println("t2 ==> ", t2)
	fmt.Println("t2 - t1 ==> ", t2-t1)
}
