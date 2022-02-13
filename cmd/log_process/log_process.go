package main

import (
	"time"

	"study/internal/log_process"
)

func main() {
	lp := log_process.NewLogProcess(
		log_process.NewReadFromFile(
			"log_file_path"),
		log_process.NewWriteToInfluxDB("influxDBDsn"),
	)
	go lp.Reader.Read(lp.Rc)
	go lp.Process()
	go lp.Writer.Write(lp.Wc)
	time.Sleep(30 * time.Second)
}
