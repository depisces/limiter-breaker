package limiter

import (
	"fmt"
	"sync"
	"time"
)

type Limiter struct {
	tb *TokenBuket
}

type TokenBuket struct {
	mu              sync.Mutex    //锁
	size            int           //桶大小
	count           int           //当前token数
	rateLimit       time.Duration //填充速率
	lastRequestTime time.Time     //最后成功请求时间
}

func (tb *TokenBuket) fillToken() {
	tb.count += tb.getTokenCount()
	if tb.count >= tb.size {
		tb.count = tb.size
	}
	fmt.Println(tb.count)
}

func (tb *TokenBuket) getTokenCount() int {
	if tb.count >= tb.size {
		return 0
	}
	if !tb.lastRequestTime.IsZero() {
		duration := time.Since(tb.lastRequestTime)
		count := int(duration / tb.rateLimit)
		return count
	}
	return 0
}

func (tb *TokenBuket) allow() bool {
	tb.fillToken()
	if tb.count > 0 {
		tb.count--
		tb.lastRequestTime = time.Now()
		return true
	}
	return false
}

// 初始化限流器
func NewLimiter(r time.Duration, size int) *Limiter {
	return &Limiter{
		tb: &TokenBuket{
			rateLimit: r,
			size:      size,
			count:     size,
		},
	}
}

func (l *Limiter) Allow() bool {
	l.tb.mu.Lock()
	defer l.tb.mu.Unlock()
	return l.tb.allow()
}
