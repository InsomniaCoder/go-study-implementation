package main

import (
	"time"
)

// Burst behavior
// ❌ Bursts are not allowed
// ❌ Excess requests are dropped or queued
// ✅ Very predictable output
type LeakyBucketRateLimiter struct {
	queue  chan struct{}
	ticker *time.Ticker
	stop   chan struct{}
}

func NewLeakyBucketRateLimiter(leakedRate time.Duration, cap int) RateLimiter {
	lb := &LeakyBucketRateLimiter{
		queue:  make(chan struct{}, cap),
		ticker: time.NewTicker(time.Second / time.Duration(leakedRate)),
		stop:   make(chan struct{}),
	}
	go lb.leak()
	return lb
}

func (rl *LeakyBucketRateLimiter) leak() {
	for {
		select {
		case <-rl.ticker.C:
			select {
			case <-rl.queue: // leaked one request
			default: // nothing to leak then move forward
			}
		case <-rl.stop:
			return
		}
	}
}

func (rl *LeakyBucketRateLimiter) Allow() bool {

	select {
	case rl.queue <- struct{}{}:
		return true
	default: //queue channel is full
		return false
	}

}
