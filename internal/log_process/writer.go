package log_process

import "fmt"

type WriteToInfluxDB struct {
	influxDBDsn string
}

func (w WriteToInfluxDB) Write(wc chan []byte) {
	for v := range wc {
		fmt.Println(string(v))
	}
}

func NewWriteToInfluxDB(influxDBDsn string) *WriteToInfluxDB {
	return &WriteToInfluxDB{influxDBDsn: influxDBDsn}
}
