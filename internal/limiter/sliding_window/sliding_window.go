package main

import (
	"fmt"
	"sync"
	"time"
)

var windowsMutex map[string]*sync.RWMutex

func init() {
	windowsMutex = make(map[string]*sync.RWMutex)
}

type timeSlot struct {
	timestamp time.Time // 这个 timeSlot 的时间起点
	count     int       // 落在这个 timeSlot 内的请求数
}

func countRequest(windows []*timeSlot) int {
	var count int
	for _, ts := range windows {
		count += ts.count
	}
	return count
}

type SlidingWindowLimiter struct {
	Slot       time.Duration // time slot 的长度
	Windows    time.Duration // sliding window 的长度
	numSlots   int           // window 内最多有多少个 slot
	maxRequest int           // window duration 内允许的最大请求数
	windows    map[string][]*timeSlot
}

func NewSlidingWindowLimiter(slot time.Duration, windows time.Duration, maxRequest int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		Slot:       slot,
		Windows:    windows,
		numSlots:   int(windows / slot),
		maxRequest: maxRequest,
		windows:    make(map[string][]*timeSlot),
	}
}

// getWindows 获取 user_id/ip 的时间窗口
func (l *SlidingWindowLimiter) getWindows(uidOrIp string) []*timeSlot {
	windows, ok := l.windows[uidOrIp]
	if !ok {
		windows = make([]*timeSlot, 0, l.numSlots)
	}
	return windows
}

func (l *SlidingWindowLimiter) storeWindow(uidOrIp string, windows []*timeSlot) {
	l.windows[uidOrIp] = windows
}

func (l *SlidingWindowLimiter) getUidOrIp() string {
	return "127.0.0.1"
}

func (l *SlidingWindowLimiter) isExceed(windows []*timeSlot) (yeah bool) {
	if countRequest(windows) < l.maxRequest {
		yeah = true
	}
	return
}

func (l *SlidingWindowLimiter) validate(uidOrIp string) (yeah bool) {
	// 同一 user_id/ip 并发安全
	mutex, ok := windowsMutex[uidOrIp]
	if !ok {
		mutex = &sync.RWMutex{}
		windowsMutex[uidOrIp] = mutex
	}
	mutex.Lock()
	defer mutex.Unlock()

	windows := l.getWindows(uidOrIp)
	now := time.Now()

	// 已经过期的 time slot 移出时间窗
	timeoutOffset := -1
	for i, ts := range windows {
		if ts.timestamp.Add(l.Windows).After(now) {
			break
		}
		timeoutOffset = i
	}
	if timeoutOffset > -1 {
		windows = windows[timeoutOffset+1:]
	}

	// 判断请求是否超限
	if l.isExceed(windows) {
		yeah = true
	}

	// 记录这次的请求数
	var lastSlot *timeSlot
	if len(windows) > 0 {
		lastSlot = windows[len(windows)-1]
		if lastSlot.timestamp.Add(l.Slot).Before(now) {
			lastSlot = &timeSlot{timestamp: now, count: 1}
			windows = append(windows, lastSlot)
		} else {
			lastSlot.count++
		}
	} else {
		lastSlot = &timeSlot{timestamp: now, count: 1}
		windows = append(windows, lastSlot)
	}

	l.storeWindow(uidOrIp, windows)

	return
}

func (l *SlidingWindowLimiter) IsLimited() bool {
	return !l.validate(l.getUidOrIp())
}

func main() {
	limiter := NewSlidingWindowLimiter(100*time.Millisecond, time.Second, 10)
	for i := 0; i < 5; i++ {
		fmt.Println(limiter.IsLimited())
	}
	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 5; i++ {
		fmt.Println(limiter.IsLimited())
	}
	fmt.Println(limiter.IsLimited())

	for _, v := range limiter.windows[limiter.getUidOrIp()] {
		fmt.Println(v.timestamp, v.count)
	}

	fmt.Println("a thousand years later...")
	time.Sleep(time.Second)
	for i := 0; i < 7; i++ {
		fmt.Println(limiter.IsLimited())
	}
	for _, v := range limiter.windows[limiter.getUidOrIp()] {
		fmt.Println(v.timestamp, v.count)
	}
}
