package main

import "fmt"

func p[T any](arr []T) {
	for _, a := range arr {
		fmt.Println(a)
	}
}

func main() {
	str := []string{"a", "b"}
	num := []int{1, 2, 3}
	p(str)
	p(num)
}
