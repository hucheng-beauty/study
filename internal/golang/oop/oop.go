package oop

import "fmt"

type Person struct {
	name string
}

// Human 多态
type Human interface {
	Speak()
}

// Chinese 封装,继承
type Chinese struct {
	p    Person
	skin string
}

func (p *Person) Walk() {
	fmt.Println(p.name + " is walk")
}

func (c *Chinese) Speak() {
	fmt.Println(c.p.name + " is speaking!")
}

type Japan struct {
	name string
	skin string
}

func (j *Japan) Speak() {
	fmt.Println(j.name + " is speaking!")
}
