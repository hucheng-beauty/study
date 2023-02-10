package pipeline

import (
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
		sli := make([]int, 0)
		for v := range in {
			sli = append(sli, v)
		}

		sort.Ints(sli)

		for _, v := range sli {
			out <- v
		}

		close(out)
	}()
	return out
}

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

func MergeN(inputs ...<-chan int) <-chan int {
	// If the length of inputs is one,return inputs[0].
	if len(inputs) > 0 && len(inputs) < 2 {
		return inputs[0]
	}

	m := len(inputs) / 2
	// Merge inputs[0:m] and inputs[m:end]
	return Merge(MergeN(inputs[:m]...), MergeN(inputs[m:]...))
}
