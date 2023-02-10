package traffic_control

import (
	"fmt"
	"testing"
)

func TestMainer(t *testing.T) {
	trafficControl := NewTrafficControl(10, 6)
	trafficCount := 100
	var aQueryCount, bQueryCount = 0, 0
	for trafficCount > 0 {
		if trafficControl.Allow() {
			aQueryCount++
		} else {
			bQueryCount++
		}
		trafficCount--
	}

	fmt.Printf("A queryCount: %v, B queryCount %v", aQueryCount, bQueryCount)
}
