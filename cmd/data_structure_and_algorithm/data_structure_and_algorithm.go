package main

import "study/internal/data_structure_and_algorithm/queue"

func main() {
	lq := queue.NewListQueue()
	lq.Enqueue("hello")
	lq.Enqueue("world")
	lq.Dequeue()
	lq.Dequeue()

}
