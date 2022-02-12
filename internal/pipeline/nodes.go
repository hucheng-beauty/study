package pipeline

import (
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"sort"
	"time"
)

/*
	pipeline:
		[ReaderSource] ==>   [node]   <==>   [node]   <==> [WriterSink]
*/

/*
	数据数据源节点,channel的关闭及检测 ==> ArraySource()
	内部排序节点 ==> InMemSort()
	归并节点 ==> Merge()
*/

var startTime time.Time

func InitTime() {
	startTime = time.Now()
}

func ArraySource(a ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

/*
	More node:
		ReaderSource node ==> ReaderSource()
		WriterSink node ==> WriterSink()
		Generate mock data ==> RandSource()
*/

// RandSource Generate rand data to <-chan int
func RandSource(count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := 0; i < count; i++ {
			out <- rand.Int()
		}
		close(out)
	}()
	return out
}

// ReaderSource Support DeBlocking
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		byteRead := 0
		for {
			n, err := reader.Read(buffer)
			byteRead += n
			if n > 0 {
				v := int(binary.BigEndian.Uint64(buffer))
				out <- v
			}
			if err != nil || (chunkSize != -1 && byteRead >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

func InMemSort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		// Read into memory
		a := []int{}
		for v := range in {
			a = append(a, v)
		}
		log.Println("Read done: ", time.Now().Sub(startTime))

		// Sort
		sort.Ints(a)
		log.Println("InMemSort done: ", time.Now().Sub(startTime))

		// Output
		for _, v := range a {
			out <- v
		}
		close(out)
	}()
	return out
}

/*
	归并排序:
		将数据分为左右俩半,分别归并排序,再把俩个有序数据排序.

	v1		   v2
	[1,3,6,7]  [1,2,3,5]   [1]
	[3,6,7]    [1,2,3,5]   [1]
	[3,6,7]    [2,3,5]     [2]
	[3,6,7]    [3,5]       [3]
	[6,7]      [3,5]       [3]
	[6,7]      [5]         [5]
	[6,7]      []          [6]
	[7]        []          [7]
*/

func Merge(in1, in2 <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		v1, ok1 := <-in1
		v2, ok2 := <-in2

		for ok1 || ok2 {
			if !ok2 || (ok1 && v1 <= v2) {
				out <- v1
				v1, ok1 = <-in1
			} else {
				out <- v2
				v2, ok2 = <-in2
			}
		}
		close(out)
		log.Println("Merge done: ", time.Now().Sub(startTime))
	}()
	return out
}

func MergeN(inputs ...<-chan int) <-chan int {
	// If the length of inputs is one,return inputs[0].
	if len(inputs) > 0 && len(inputs) < 2 {
		return inputs[0]
	}

	m := len(inputs) / 2
	// Merge inputs[0:m] and inputs[m:end]
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}

func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}
