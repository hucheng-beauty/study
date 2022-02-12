package policy_mode

import "fmt"

/*
	对象行为型模式:
		策略模式: 作为一种软件设计模式,指对象有某个行为,但是在不同的场景中,该行为有不同的实现算法
		eg. 个人所得税:	美国个税与中国个税有着不同的算法

					 依赖 				    实现
	context上下文类  ==========>  抽象策略类 =========> 具体的策略1
													  具体的策略2
					  调用
	context上下文类 ==========>	具体的策略

	优点:
		完美支持"开闭原则"
		避免适用多重条件转移语句
		提供类管理相关的算法族的办法

	缺点:
		客户端必须知道所以的策略类,并自己决定使用哪一个策略类
		策略模式将会造成产生很多策略类

	适用场景:
		需要动态的在几种算法中选择一种
		多个类区别仅在于他们的行为或算法不同的场景
*/

// Context 实现一个 context 上下文类
type Context struct {
	Strategy
}

// Strategy 抽象的策略
type Strategy interface {
	Do()
}

// Strategy1 实现具体的策略: 策略1
type Strategy1 struct{}

func (s *Strategy1) Do() {
	fmt.Println("Strategy1")
}

// Strategy2 实现具体的策略: 策略2
type Strategy2 struct{}

func (s *Strategy2) Do() {
	fmt.Println("Strategy2")
}
