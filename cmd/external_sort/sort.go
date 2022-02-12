package main

import (
	"bufio"
	"os"
	"strconv"

	"study/pkg/pipeline"
)

func createNetworkPipeline(filename string, fileSize, chunkCount int) <-chan int {
	sortAddr := []string{}
	pipeline.InitTime()
	for i := 0; i < chunkCount; i++ {
		chunkSize := fileSize / chunkCount // TODO:  need to deal with chunkCount is not divisible
		file, err := os.Open(filename)     // TODO: file is not close, need to deal with
		if err != nil {
			panic(err)
		}

		file.Seek(int64(i*chunkSize), 0)

		source := pipeline.ReaderSource(bufio.NewReader(file), chunkSize)

		addr := ":" + strconv.Itoa(7000+i)
		pipeline.NetworkSink(addr, pipeline.InMemSort(source))

		sortAddr = append(sortAddr, addr)
	}

	sortResults := []<-chan int{}
	for _, addr := range sortAddr {
		sortResults = append(sortResults, pipeline.NetworkSource(addr))
	}
	return pipeline.MergeN(sortResults...)
}
