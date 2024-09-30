package agency_mode

import (
	"fmt"
)

/*
   代理模式
        代理模式是一种结构型设计模式，它提供了一个代理对象，用于控制对目标对象的访问。
        代理对象可以对目标对象进行预处理、过滤、后处理等操作，并将其结果返回给客户端，从而实现对目标对象的访问控制。
*/

type PaymentService interface {
	pay(order string) string
}

type WXPay struct{}

func (w *WXPay) pay(order string) string {
	return "从微信获取支付token" + order
}

type AliPay struct{}

func (a *AliPay) pay(order string) string {
	return "从阿里获取支付token" + order
}

// PaymentProxy 代理对象
type PaymentProxy struct {
	realPay PaymentService
}

func NewPaymentProxy(realPay PaymentService) *PaymentProxy {
	return &PaymentProxy{realPay: realPay}
}

// Pay 实现预处理、过滤、后处理等操作
func (p *PaymentProxy) Pay(order string) string {
	// 功能增强
	fmt.Println("处理" + order)
	fmt.Println("1、校验签名")
	fmt.Println("2、格式化订单数据")
	fmt.Println("3、参数检查")
	fmt.Println("4、记录请求日志")

	token := p.realPay.pay(order)

	url := "http://组装" + token + "然后跳转到第三方支付"
	return url
}
