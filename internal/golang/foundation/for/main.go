package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(file)
	// 省略初始值和递增表达式
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func forever() {
	// 省略初始值,结束条件,递增表达式
	for {
		fmt.Println("I love you.")
	}
}

func convertToBin(n int) string {
	result := ""
	// 省略初始值
	for ; n > 0; n /= 2 {
		result = strconv.Itoa(n%2) + result
	}
	return result
}

func main() {
	fmt.Println(
		convertToBin(10),
		convertToBin(0),
		convertToBin(8),
		convertToBin(2121),
	)

	//forever()
	printFile("./internal/golang/foundation/for/abc.txt")
}
