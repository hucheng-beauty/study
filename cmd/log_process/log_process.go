package main

import (
	"time"

	"study/internal/log_process"
)

func main() {
	filePath := "log_file_path"
	influxDBDSN := "influxDBDSN"

	lp := log_process.NewLogProcess(
		log_process.NewReadFromFile(filePath),
		log_process.NewWriteToInfluxDB(influxDBDSN),
	)
	go lp.Reader.Read(lp.Rc)
	go lp.Process()
	go lp.Writer.Write(lp.Wc)
	time.Sleep(30 * time.Second)
}
