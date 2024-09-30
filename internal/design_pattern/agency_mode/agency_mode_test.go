package agency_mode

import (
	"fmt"
	"testing"
)

func TestMainer(t *testing.T) {
	paymentProxy := NewPaymentProxy(&AliPay{})
	url := paymentProxy.Pay("阿里订单")

	fmt.Println(url)
}
