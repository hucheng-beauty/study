package main

import "strings"

func Worker(str string) (string, error) {
	return strings.ToUpper(str), nil
}
