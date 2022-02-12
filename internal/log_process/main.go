package main

/*
	1. 从文件中读取日志
	2. 解析日志
	3. 写入数据库
*/

type LogProcess struct {
	path string
}

func (lp *LogProcess) ReadFromFile() {
	// 读取日志
}

func (lp *LogProcess) Process() {
	// 解析日志
}

func (lp *LogProcess) WriteToInfluxDB() {
	// 写入数据库
}

func main() {

}
