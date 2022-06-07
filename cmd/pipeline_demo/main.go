package main

import (
	"bufio"
	"fmt"
	"os"

	"study/internal/pipeline"
)

func main() {
	const filename = "large.in"
	const n = 100000000

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.RandSource(n)
	writer := bufio.NewWriter(file)
	pipeline.WriterSink(writer, p)
	writer.Flush()

	file, err = os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p = pipeline.ReaderSource(bufio.NewReader(file), 1)
	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}

func mergeDemo() {
	p := pipeline.Merge(
		pipeline.InMemorySort(pipeline.ArraySource(3, 1, 9, 5, 7)),
		pipeline.InMemorySort(pipeline.ArraySource(4, 2, 8, 6, 10)))
	for v := range p {
		fmt.Printf("%d ", v)
	}
}
