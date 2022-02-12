package pipeline

import (
	"bufio"
	"net"
)

// NetworkSink Write the data in the network into the <-chan int
func NetworkSink(addr string, in <-chan int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		defer listener.Close()

		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		// use buffer area
		writer := bufio.NewWriter(conn)
		defer writer.Flush()

		WriterSink(writer, in)
	}()
}

// NetworkSource Read data form network.
func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		// use buffer area to speed up  the reading
		reader := bufio.NewReader(conn)
		r := ReaderSource(reader, -1)
		for v := range r {
			out <- v
		}

		close(out)
	}()
	return out
}
