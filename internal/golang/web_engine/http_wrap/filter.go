package main

import (
	"log"
	"time"
)

type FilterBuilder func(next Filter) Filter

type Filter handleFunc

func MetricFilterBuilder(next Filter) Filter {
	return func(ctx *Context) {
		start := time.Now().Nanosecond()
		next(ctx)
		end := time.Now().Nanosecond()
		log.Printf("used %d nanosecond", end-start)
	}
}
