package main

import (
	"bufio"
	"fmt"
	"os"
)

func writeFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_EXCL|os.O_CREATE, 0666)
	if err != nil {
		if pathError, ok := err.(*os.PathError); ok {
			fmt.Println(pathError.Op, pathError.Path, pathError.Error())
		} else {
			panic(err)
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
