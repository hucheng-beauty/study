package main

import (
	"encoding/json"
	"fmt"
)

/*
	练习:自然语言处理小工具
*/

/*
	json的解析:
		1.json 数据格式
		2.结构体 tag
		3.json 的 Marshal & UnMarshal,数据类型:
			Marshal(struct => json string)
			UnMarshal(json string => struct)
		4.第三方API的解析技巧
*/

type Order struct {
	Id         string  `json:"id"`
	Name       string  `json:"name,omitempty"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}

func parseNPL() {

}

func main() {
	o := Order{
		Id:         "1234",
		Name:       "learn go",
		Quantity:   3,
		TotalPrice: 30,
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", bytes)
}
