package main

import (
	"sync"
	"time"
)

// good for handling burst
// there's a bucket and each request burns it
// then it generates back for each second pass
// the rate per second then is used to control (smoothen)
// while cap is the upper worst case that the system can handle per second
type TokenBucketRateLimiter struct {
	mu            sync.Mutex
	rps           int
	cap           int
	bucket        int
	latestRequest time.Time
}

func NewTokenBucketRateLimiter(rps int, cap int) RateLimiter {
	return &TokenBucketRateLimiter{
		rps:           rps,
		cap:           cap,
		bucket:        cap,
		latestRequest: time.Now(),
	}
}

func (rl *TokenBucketRateLimiter) Allow() bool {

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	timePass := now.Sub(rl.latestRequest).Seconds()
	// This calculates fractional tokens first, then converts to int when adding to the bucket. For example, with 10 rps and 0.5 seconds passed, you'd get 5 tokens instead of 0.
	tokenToAdd := timePass * float64(rl.rps)
	rl.bucket = min(rl.cap, rl.bucket+int(tokenToAdd))
	rl.latestRequest = now

	if rl.bucket > 0 {
		rl.bucket--
		return true
	}
	return false
}
