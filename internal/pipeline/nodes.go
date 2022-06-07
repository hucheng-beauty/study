package pipeline

import (
	"encoding/binary"
	"io"
	"math/rand"
	"sort"
	"time"
)

/*
	pipeline:
		[ReaderSource] ==>   [node]   <==>   [node]   <==> [WriterSink]
*/

var startTime time.Time

func InitTime() {
	startTime = time.Now()
}

// ArraySource Put the int data of array to <-chan int
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

// InMemorySort sort data from <-chan int
// 归并排序:将数据分为左右俩半,分别归并排序,再把俩个有序数据排序
// v1		   v2
// [1,3,6,7]  [1,2,3,5]   [1]
// [3,6,7]    [1,2,3,5]   [1]
// [3,6,7]    [2,3,5]     [2]
// [3,6,7]    [3,5]       [3]
// [6,7]      [3,5]       [3]
// [6,7]      [5]         [5]
// [6,7]      []          [6]
// [7]        []          [7]
func InMemorySort(in <-chan int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		// Read into memory
		memSli := make([]int, 0)
		for v := range in {
			memSli = append(memSli, v)
		}

		// Sort
		sort.Ints(memSli)

		// Output
		for _, v := range memSli {
			out <- v
		}

		close(out)
	}()
	return out
}

// Merge merge in1 in2 and put into <-chan
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
	}()
	return out
}

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

// WriterSink put data from in and write into writer
func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
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
