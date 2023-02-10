package response

import "reflect"

type Response struct {
	Action  string      `json:"Action"`
	RtCode  int         `json:"RtCode"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data"`
}

func New(action string, rtCode int, message string, data interface{}) *Response {
	return &Response{
		Action:  action + "Response",
		RtCode:  rtCode,
		Message: message,
		Data:    data,
	}
}

type Pagination struct {
	TotalCount int `json:"TotalCount"`
	Offset     int `json:"Offset"`
	Limit      int `json:"Limit"`
}

func NewPagination(totalCount int, offset int, limit int) *Pagination {
	return &Pagination{
		TotalCount: totalCount,
		Offset:     offset,
		Limit:      limit,
	}
}

type WithPagination struct {
	*Pagination
	*Response
}

func NewWithPagination(response *Response, pagination *Pagination) *WithPagination {
	return &WithPagination{
		Response:   response,
		Pagination: pagination,
	}
}

func Success(action string, data interface{}) *Response {
	return New(action, 0, "", data)
}

func SuccessWithPagination(action string, data interface{}, pagination *Pagination) *WithPagination {
	if v := reflect.ValueOf(data); v.IsNil() {
		data = []interface{}{}
	}
	response := New(action, 0, "", data)
	return NewWithPagination(response, pagination)
}

func Error(action string, rtCode int, message string) *Response {
	return New(action, rtCode, message, nil)
}

/*
	TODO add response list
*/
