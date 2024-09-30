package main

import "fmt"

func main() {
	// arr := [...]int{0, 1, 2, 3, 4, 5, 6, 7}
	//
	// // slice 不可向前扩展,可以向后扩展
	//
	// s1 := arr[2:6] // s1 = [2,3,4,5], len(s1) = 4, cap(s1) = 6
	//
	// // s[i] 不可超越 len(s);向后扩展不可超越底层数组 cap(s)
	// s2 := s1[3:5] // s2 = [5, 6], len(s2) = 2, cap(s2) = 3
	//
	// fmt.Printf("s1=%v, len(s1)=%d, cap(s1)=%d\n", s1, len(s1), cap(s1))
	// fmt.Printf("s2=%v, len(s2)=%d, cap(s2)=%d\n", s2, len(s2), cap(s2))
	a := 1
	b := 2

	sli := []*int{&a, &b}
	for _, i := range sli {
		fmt.Printf("%d\n", *i)
	}
	fmt.Printf("sli=%v, len(sli)=%d, cap(sli)=%d\n", sli, len(sli), cap(sli))

	T(sli)
	for _, i := range sli {
		fmt.Printf("%d\n", *i)
	}
	fmt.Printf("sli=%v, len(sli)=%d, cap(sli)=%d\n", sli, len(sli), cap(sli))
}

func T(sli []*int) {
	c := 3
	sli[0] = &c

	d := 4
	sli = append(sli, &d)
}
