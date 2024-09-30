package main

import (
	"fmt"
	"strconv"
)

/*
	定义
		使多个对象都有机会处理请求，从而避免请求的发送者和接收者之间的耦合关系，
		将这个对象连成一条链，并沿着这条链传递该请求，直到有一个对象处理它为止

	优点
		降低了对象之间的耦合度。该模式使得一个对象无须知道到底是哪一个对象处理其请求以及链的结构，发送者和接收者也无须拥有对方的明确信息。
		增强了系统的可扩展性。可以根据需要增加新的请求处理类，满足开闭原则。
		增强了给对象指派职责的灵活性。当工作流程发生变化，可以动态地改变链内的成员或者调动它们的次序，也可动态地新增或者删除责任。
		责任链简化了对象之间的连接。每个对象只需保持一个指向其后继者的引用，
			不需保持其他所有处理者的引用，这避免了使用众多的 if 或者 if···else 语句。
		责任分担。每个类只需要处理自己该处理的工作，不该处理的传递给下一个对象完成，明确各类的责任范围，符合类的单一职责原则。
	缺点
		不能保证每个请求一定被处理。由于一个请求没有明确的接收者，所以不能保证它一定会被处理，该请求可能一直传到链的末端都得不到处理。
		对比较长的职责链，请求的处理可能涉及多个处理对象，系统性能将受到一定影响。
		职责链建立的合理性要靠客户端来保证，增加了客户端的复杂性，可能会由于职责链的错误设置而导致系统出错，如可能会造成循环调用。
		场景
		有请求需要处理，但是不知道谁处理时

*/

type Handler interface {
	Handler(handlerID int) string
}

type handler struct {
	name         string
	next         Handler
	RequestLevel int
}

func NewHandler(name string, next Handler, requestLevel int) *handler {
	return &handler{
		name:         name,
		next:         next,
		RequestLevel: requestLevel,
	}
}

func (h *handler) Handler(requestLevel int) string {
	if h.RequestLevel == requestLevel {
		return h.name + " handled " + strconv.Itoa(requestLevel)
	}
	if h.next == nil {
		return ""
	}
	return h.next.Handler(requestLevel)
}

func main() {
	wang := NewHandler("wang", nil, 1)
	zhang := NewHandler("zhang", wang, 2)
	wu := NewHandler("wu", zhang, 3)
	r := wu.Handler(2)
	fmt.Println(r)
}
