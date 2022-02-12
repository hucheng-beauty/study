package main

import (
	"bufio"
	"fmt"
	"os"

	"study/internal/pipeline"
)

func main() {
	// p := createPipeline("large.in", 800*1000*100, 4)
	// writeToFile(p, "large.out")
	// printFile("large.out")

	p := createNetworkPipeline("large.in", 800*1000*1000, 4)
	writeToFile(p, "large.out")
	printFile("large.out")
}

func createPipeline(filename string, fileSize, chunkCount int) <-chan int {
	sortResults := []<-chan int{}
	pipeline.InitTime()
	for i := 0; i < chunkCount; i++ {
		chunkSize := fileSize / chunkCount // TODO:  need to deal with chunkCount is not divisible
		file, err := os.Open(filename)     // TODO: file is not close, need to deal with
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)

		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)
		sortResults = append(sortResults, pipeline.InMemSort(source))
	}

	return pipeline.MergeN(sortResults...)
}

func writeToFile(p <-chan int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	pipeline.WriterSink(writer, p)
}

func printFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	p := pipeline.ReaderSource(file, -1)

	count := 0
	for v := range p {
		fmt.Println(v)
		count++
		if count >= 100 {
			break
		}
	}
}
