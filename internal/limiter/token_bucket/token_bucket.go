package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket 令牌桶限流
type TokenBucket struct {
	mutex     sync.Mutex
	capacity  int           // 桶的容量
	rate      time.Duration // 多长时间生成一个令牌
	tokens    int           // 桶中当前 token 数量
	lastToken time.Time     // 桶上次放 token 的时间戳 s
}

func NewTokenBucket(bucketSize int, rate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity: bucketSize,
		rate:     rate,
	}
}

// Allow 验证是否能获取一个令牌 返回是否被限流
func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()

	if tb.lastToken.IsZero() {
		// 第一次访问初始化为最大令牌数
		tb.lastToken, tb.tokens = now, tb.capacity
	} else {
		if tb.lastToken.Add(tb.rate).After(now) {
			// 如果 now 与上次请求的间隔超过了 token rate
			// 则增加令牌,更新 lastToken
			tb.tokens += int(now.Sub(tb.lastToken) / tb.rate)
			if tb.tokens > tb.capacity {
				tb.tokens = tb.capacity
			}
			tb.lastToken = now
		}
	}

	if tb.tokens > 0 {
		// 如果令牌数大于0，取走一个令牌
		tb.tokens--
		return true
	}

	// 没有令牌,则拒绝
	return false
}

func main() {
	tokenBucket := NewTokenBucket(5, 100*time.Millisecond)
	for i := 0; i < 10; i++ {
		fmt.Println(tokenBucket.Allow())
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println(tokenBucket.Allow())
}
