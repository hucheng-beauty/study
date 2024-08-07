package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶限流
type TokenBucket struct {
	mutex     sync.Mutex // 互斥锁
	capacity  int64      // 桶的容量
	rate      float64    // 令牌放入速率
	tokens    float64    // 当前令牌数量
	lastToken time.Time  // 上一次放令牌的时间
}

func (tb *TokenBucket) Allow() bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	now := time.Now()

	// 计算需要放的令牌数量
	tb.tokens = tb.tokens + tb.rate*now.Sub(tb.lastToken).Seconds()
	if tb.tokens > float64(tb.capacity) {
		tb.tokens = float64(tb.capacity)
	}

	// 判断是否允许请求
	if tb.tokens >= 1 {
		tb.tokens--
		tb.lastToken = now
		return true
	}

	return false
}

// LimitHandler 中间件使用限流
func LimitHandler(maxConn int /*最大连接数*/) gin.HandlerFunc {
	tb := &TokenBucket{
		capacity:  int64(maxConn),
		rate:      1.0,
		tokens:    0,
		lastToken: time.Now(),
	}
	return func(c *gin.Context) {
		if !tb.Allow() {
			c.String(http.StatusServiceUnavailable /*503*/, "Too many request")
			c.Abort()
			return
		}
		c.Next()
	}
}
