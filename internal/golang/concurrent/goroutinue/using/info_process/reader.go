package main

import "sync/atomic"

type ReadFromRemote struct{}

var count int64 = 0

func (ReadFromRemote) Read(rc chan string) {
	for i := 0; i < 100; i++ {
		rc <- "golang"
		atomic.AddInt64(&count, 1)
	}
}

func NewReadFromRemote() *ReadFromRemote {
	return &ReadFromRemote{}
}
