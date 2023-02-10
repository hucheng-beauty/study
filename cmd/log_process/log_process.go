package main

import (
	"time"

	"study/internal/log_process"
)

const (
	filePath    = "log_file_path"
	influxDBDSN = "influxDBDSN"
)

func main() {
	lp := log_process.NewLogProcess(
		log_process.NewReadFromFile(filePath),
		log_process.NewWriteToInfluxDB(influxDBDSN),
	)

	go lp.Reader.Read(lp.Rc)
	go lp.Process()
	go lp.Writer.Write(lp.Wc)

	time.Sleep(30 * time.Second)
}
