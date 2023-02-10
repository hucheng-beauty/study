package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func writeFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_EXCL|os.O_CREATE, 0666)
	if err != nil {
		// error asserting
		var pathError *os.PathError
		if errors.As(err, &pathError) {
			fmt.Println(pathError.Op, pathError.Path, pathError.Error())
		}
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < 20; i++ {
		fmt.Fprintln(writer, "hello")
	}
}

func main() {
	writeFile("test.txt")
}
