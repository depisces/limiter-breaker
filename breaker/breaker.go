package breaker

import (
	"errors"
	"sync"
	"time"
)

const (
	STATE_CLOSE     = iota //限流器关闭
	STATE_OPEN             //限流器打开
	STATE_HALF_OPEN        //限流器半开
)

type Breaker struct {
	mu                sync.Mutex
	state             int
	failureThreshold  int
	successThreshold  int
	halfMaxRequest    int
	halfCycleReqCount int
	timeout           time.Duration
	failureCount      int
	successCount      int
	cycleStartTime    time.Time
}

func NewBreaker(failureThreshold, successThreshold, halfMaxRequest int, timeout time.Duration) *Breaker {
	return &Breaker{
		state:            STATE_CLOSE,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		halfMaxRequest:   halfMaxRequest,
		timeout:          timeout,
	}
}

func (b *Breaker) Exec(f func() error) error {
	b.before()
	if b.state == STATE_OPEN {
		return errors.New("断路器打开")
	}
	if b.state == STATE_CLOSE {
		err := f()
		b.after(err)
		return err
	}
	if b.state == STATE_HALF_OPEN {
		if b.halfCycleReqCount < b.halfMaxRequest {
			err := f()
			b.after(err)
			return err
		} else {
			return errors.New("短时间请求次数过多")
		}
	}
	return nil
}

func (b *Breaker) before() {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case STATE_OPEN:
		if b.cycleStartTime.Add(b.timeout).Before(time.Now()) {
			b.state = STATE_HALF_OPEN
			b.reset()
			return
		}
	case STATE_HALF_OPEN:
		if b.successCount >= b.successThreshold {
			b.state = STATE_CLOSE
			b.reset()
			return
		}
		if b.cycleStartTime.Add(b.timeout).Before(time.Now()) {
			b.cycleStartTime = time.Now()
			b.halfCycleReqCount = 0
			return
		}
	case STATE_CLOSE:
		if b.cycleStartTime.Add(b.timeout).Before(time.Now()) {
			b.reset()
		}
	}

}

func (b *Breaker) after(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err == nil {
		b.onSuccess()
	} else {
		b.onFailure()
	}
}

func (b *Breaker) onSuccess() {
	b.failureCount = 0
	if b.state == STATE_HALF_OPEN {
		b.successCount++
		b.halfCycleReqCount++
		if b.successCount >= b.successThreshold {
			b.state = STATE_CLOSE
			b.reset()
		}
	}
}
func (b *Breaker) onFailure() {
	b.successCount = 0
	b.failureCount++
	if b.state == STATE_HALF_OPEN || (b.state == STATE_CLOSE && b.failureCount >= b.failureThreshold) {
		b.state = STATE_OPEN
		b.reset()
		return
	}

}

func (b *Breaker) reset() {
	b.successCount = 0
	b.failureCount = 0
	b.halfCycleReqCount = 0
	b.cycleStartTime = time.Now()
}
