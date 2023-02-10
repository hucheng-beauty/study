package main

import (
	"fmt"
)

type Message struct {
	msg string
}

func NewMessage(msg string) Message {
	return Message{
		msg: msg,
	}
}

type Greeter struct {
	Message Message
}

func NewGreeter(m Message) Greeter {
	return Greeter{Message: m}
}

type Event struct {
	Greeter Greeter
}

func NewEvent(g Greeter) Event {
	return Event{Greeter: g}
}
func (e Event) Start() {
	fmt.Println(e.Greeter.Greet())
}
func (g Greeter) Greet() Message {
	return g.Message
}

// 使用wire前

// func main() {
//	message := NewMessage("hello world")
//	greeter := NewGreeter(message)
//	event := NewEvent(greeter)
//
//	event.Start()
// }

// 使用wire后
func main() {
	event := InitializeEvent("hello world")
	event.Start()
}
