package main

import "time"

type Person struct {
	Name     string `form:"name"`
	Address  string `form:"address"`
	Birthday string `form:"birthday" time_format:"2006-01-02" `
}

type People struct {
	Age     int    `form:"age" binding:"required,gt=10"`
	Name    string `form:"name" binding:"required"`
	Address string `form:"address" binding:"required"`
}

type Booking struct {
	CheckIn  time.Time `form:"check_in" validate:"required, bookableDate" time_format:"2006-01-02"`
	CheckOut time.Time `form:"check_out" validate:"required, gtfield=CheckIn" time_format:"2006-01-02"`
}
