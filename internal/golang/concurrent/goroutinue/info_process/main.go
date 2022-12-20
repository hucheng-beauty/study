package main

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}

type Process struct {
	Rc     chan string
	Wc     chan string
	Reader Reader
	Writer Writer
}

func (app *Process) Process() {
	for v := range app.Rc {
		app.Wc <- strings.ToUpper(v)
	}
}

func NewProcess(reader Reader, writer Writer) *Process {
	return &Process{
		Rc:     make(chan string),
		Wc:     make(chan string),
		Reader: reader,
		Writer: writer,
	}
}

func main() {
	app := NewProcess(NewReadFromRemote(), NewWriteToMongoDB())

	for i := 0; i < 1000; i++ {
		go app.Reader.Read(app.Rc)
	}

	go app.Process()
	for i := 0; i < 10; i++ {
		go app.Writer.Write(app.Wc)
	}

	time.Sleep(2 * time.Millisecond)
	fmt.Println(atomic.LoadInt64(&count))

	panic("stack")
}
