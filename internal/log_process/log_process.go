package log_process

import (
	"strings"
)

/*
	1. 从文件中读取日志
	2. 解析日志
	3. 写入数据库
*/

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan []byte)
}

type LogProcess struct {
	Rc     chan []byte
	Wc     chan []byte
	Reader Reader
	Writer Writer
}

func (lp *LogProcess) Process() {
	/*
		1. 从 Read Channel 中读取每行日志数据
		2. 正则提取所需要的监控数据(path、status、method等)
		3. 写入 Write Channel
	*/
	for v := range lp.Rc {
		lp.Wc <- []byte(strings.ToUpper(string(v)))
	}
}

func NewLogProcess(read Reader, write *WriteToInfluxDB) *LogProcess {
	return &LogProcess{
		Rc:     make(chan []byte),
		Wc:     make(chan []byte),
		Reader: read,
		Writer: write,
	}
}
