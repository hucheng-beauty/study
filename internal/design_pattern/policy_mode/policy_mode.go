package policy_mode

import "fmt"

/*
	对象行为型模式:
		策略模式: 作为一种软件设计模式,指对象有某个行为,但是在不同的场景中,该行为有不同的实现算法
		eg. 个人所得税:	美国个税与中国个税有着不同的算法

			       依赖 		    实现
	LoggerManager  ==>  Logger  ==>  DB

			       调用
	LoggerManager  ==>  File/DB

	优点:
		完美支持"开闭原则"
		避免适用多重条件转移语句
		提供类管理相关的算法族的办法

	缺点:
		客户端必须知道所以的策略类,并自己决定使用哪一个策略类
		策略模式将会造成产生很多策略类

	适用场景:
		需要动态的在几种算法中选择一种
		多个类区别仅在于他们的行为或算法不同的场景
*/

// LoggerManager 实现一个日志记录器(相当于Context)
type LoggerManager struct {
	Logger
}

func NewLoggerManager(logger Logger) *LoggerManager {
	return &LoggerManager{Logger: logger}
}

// Logger 抽象的日志
type Logger interface {
	Info()
	Error()
}

// FileLog 实现具体的日志: 文件方式记录日志
type FileLog struct{}

func NewFileLog() *FileLog {
	return &FileLog{}
}

func (FileLog) Info() {
	fmt.Println("FileLog: Info")
}

func (FileLog) Error() {
	fmt.Println("FileLog: Error")
}

// DBLog 实现具体的日志: 数据库方式记录日志
type DBLog struct{}

func NewDBLog() *DBLog {
	return &DBLog{}
}

func (DBLog) Info() {
	fmt.Println("DBLog: Info")
}

func (DBLog) Error() {
	fmt.Println("DBLog: Error")
}
