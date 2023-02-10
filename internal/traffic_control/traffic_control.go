package traffic_control

import (
	"math/rand"
	"sync/atomic"
)

// TrafficControl balanced traffic segmentation.
type TrafficControl struct {
	source     []int
	queryCount uint32
	base       int // 基准
	ratio      int // 比例
}

func NewTrafficControl(base int, ratio int) *TrafficControl {
	source := make([]int, base)
	// source = [0 1 2 3 4 5 6 7 8 9]
	for i := 0; i < base; i++ {
		source[i] = i
	}

	// 将 source 中的数据打乱
	// source = [1 7 4 0 9 2 3 5 8 6]
	rand.Shuffle(base, func(i, j int) {
		source[i], source[j] = source[j], source[i]
	})

	return &TrafficControl{
		source: source,
		base:   base,
		ratio:  ratio,
	}
}

func (t *TrafficControl) Allow() bool {
	rate := t.source[int(atomic.AddUint32(&t.queryCount, 1))%t.base]
	if rate < t.ratio {
		return true
	} else {
		return false
	}
}
