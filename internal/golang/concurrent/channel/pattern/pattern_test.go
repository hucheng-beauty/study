package main

import (
	"testing"
)

func TestMainer(t *testing.T) {
	s1 := []string{"c", "cpp", "go", "python", "java", "rust", "ruby"}
	ch1 := make(chan string, 1)
	s2 := []string{"A", "B", "C", "D", "E", "F", "G"}
	ch2 := make(chan string, 1)
	s3 := []string{"a", "b", "c", "d", "e", "f", "g"}
	ch3 := make(chan string, 1)

	for _, s := range s1 {
		ch1 <- s
	}

	for _, s := range s2 {
		ch1 <- s
	}

	for _, s := range s3 {
		ch1 <- s
	}

	fanInNew(ch1, ch2, ch3)
}
