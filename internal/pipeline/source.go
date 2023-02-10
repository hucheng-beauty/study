package pipeline

import (
	"bufio"
	"encoding/binary"
	"io"
	"math/rand"
	"net"
)

// ArraySource generate array data.
func ArraySource(in ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, v := range in {
			out <- v
		}
		close(out)
	}()
	return out
}

// RandSource generate rand data.
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

// ReaderSource support DeBlocking
func ReaderSource(reader io.Reader, chunkSize int) <-chan int {
	out := make(chan int, 1024)
	go func() {
		buffer := make([]byte, 8)
		readCount := 0
		for {
			n, err := reader.Read(buffer)

			readCount += n
			if n > 0 {
				out <- int(binary.BigEndian.Uint64(buffer))
			}

			if err != nil || (chunkSize != -1 && readCount >= chunkSize) {
				break
			}
		}
		close(out)
	}()
	return out
}

// NetworkSource read data form network.
func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		// use buffer reader to speed up the reading
		reader := bufio.NewReader(conn)
		r := ReaderSource(reader, -1)
		for v := range r {
			out <- v
		}

		close(out)
	}()
	return out
}
