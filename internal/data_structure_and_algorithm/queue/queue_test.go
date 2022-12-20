package queue

import (
	"testing"
)

func TestMainer(t *testing.T) {
	lq := NewListQueue()
	lq.Enqueue("hello")
	lq.Enqueue("world")
	lq.Dequeue()
	lq.Dequeue()
}
