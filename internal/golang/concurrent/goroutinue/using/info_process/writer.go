package main

import "fmt"

type WriteToMongoDB struct{}

func (WriteToMongoDB) Write(wc chan string) {
	for v := range wc {
		fmt.Println(v) // insert into mongo db
	}
}

func NewWriteToMongoDB() *WriteToMongoDB {
	return &WriteToMongoDB{}
}
