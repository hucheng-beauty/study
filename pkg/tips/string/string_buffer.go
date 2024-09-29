package main

import (
	"bytes"
	"fmt"
	"time"
)

// StringBuffer buffer 字符串拼接
func StringBuffer() {
	t1 := time.Now().UnixNano() / 1e6
	fmt.Println("t1 ==> ", t1)

	s := bytes.NewBufferString("hello")
	for i := 0; i < 500000; i++ {
		s.WriteString("world")
	}

	t2 := time.Now().UnixNano() / 1e6
	fmt.Println("t2 ==> ", t2)
	fmt.Println("t2 - t1 ==> ", t2-t1)
}
