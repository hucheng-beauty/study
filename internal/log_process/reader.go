package log_process

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type ReadFromFile struct {
	path string // 读取文件路径
}

func (r ReadFromFile) Read(rc chan []byte) {
	/*
		1. 打开文件
		2. 从末尾逐行读取日志内容
	*/

	// open the file
	file, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error: %s", err.Error()))
	}

	// read content from the end of file
	file.Seek(0, 2)
	rd := bufio.NewReader(file)
	for {
		line, errRead := rd.ReadBytes('\n')
		if errRead == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			fmt.Println(fmt.Sprintf("ReadBytes error: %s", err.Error()))
			continue
		}
		rc <- line[:len(line)-1] // remove '\n'
	}
}

func NewReadFromFile(path string) *ReadFromFile {
	return &ReadFromFile{path: path}
}
