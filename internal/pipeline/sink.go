package pipeline

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

// NetworkSink write network data into the <-chan int
func NetworkSink(addr string, in <-chan int) {
	listener, errListen := net.Listen("tcp", addr)
	if errListen != nil {
		return
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

// WriterSink write data into writer.
func WriterSink(writer io.Writer, in <-chan int) {
	for v := range in {
		buffer := make([]byte, 8)
		binary.BigEndian.PutUint64(buffer, uint64(v))
		writer.Write(buffer)
	}
}
