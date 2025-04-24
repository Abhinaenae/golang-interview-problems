package main

import "time"

type RateLimiter struct {
	ticker *time.Ticker
}

func NewRateLimiter(n int) *RateLimiter {
	opsPerSecond := time.Duration(n)
	limit := time.Second / opsPerSecond
	return &RateLimiter{
		ticker: time.NewTicker(limit),
	}
}

func (r *RateLimiter) CanTake() bool {
	select {
	case <-r.ticker.C:
		return true
	default:
		return false
	}
}

func (r *RateLimiter) Take() {
	<-r.ticker.C
}
