package example

import (
	"fmt"
	"sort"
)

func ExampleMax() {
	fmt.Println(Max(1, 2))
	// Output:
	// 2
}

func Example() {
	people := []Person{
		{"Bob", 31},
		{"John", 42},
		{"Michael", 17},
		{"Jenny", 26},
	}

	fmt.Println(people)
	sort.Sort(ByAge(people))
	fmt.Println(people)

	// Output:
	// [Bob: 31 John: 42 Michael: 17 Jenny: 26]
	// [Michael: 17 Jenny: 26 Bob: 31 John: 42]
}
