package limit

import (
	"sync/atomic"
	"time"
)

type Limiter struct {
	Interval time.Duration
	MaxCount int32
	ReqCount int32
}

func NewLimiter(interval time.Duration, maxNum int32) *Limiter {
	Limit := &Limiter{
		Interval: interval,
		MaxCount: maxNum,
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			<-ticker.C
			Limit.set(-Limit.Get())
		}
	}()

	return Limit
}

func (l *Limiter) Increase() bool {
	if l.Get() < l.MaxCount {
		l.set(1)
		return true
	}
	return false
}

func (l *Limiter) set(v int32) {
	atomic.AddInt32(&l.ReqCount, v)
}

func (l *Limiter) Get() int32 {
	return atomic.LoadInt32(&l.ReqCount)
}
