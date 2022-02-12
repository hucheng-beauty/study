package policy_mode

import "fmt"

// 实现一个日志记录器(相当于Context)

type LoggerManager struct {
	Logger
}

func NewLoggerManager(logger Logger) *LoggerManager {
	return &LoggerManager{Logger: logger}
}

// 抽象的日志

type Logger interface {
	Info()
	Error()
}

// 实现具体的日志: 文件方式记录日志

type FileLog struct{}

func (FileLog) Info() {
	fmt.Println("FileLog: Info")
}

func (FileLog) Error() {
	fmt.Println("FileLog: Error")
}

// 实现具体的日志: 数据库方式记录日志

type DBLog struct{}

func (DBLog) Info() {
	fmt.Println("DBLog: Info")
}

func (DBLog) Error() {
	fmt.Println("DBLog: Error")
}

func main() {
	fileLog := &FileLog{}
	loggerManage := NewLoggerManager(fileLog)

	loggerManage.Info()
	loggerManage.Error()

	dbLog := &DBLog{}
	loggerManage.Logger = dbLog

	loggerManage.Info()
	loggerManage.Error()

}
