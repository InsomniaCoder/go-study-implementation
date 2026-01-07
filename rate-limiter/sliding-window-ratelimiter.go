package main

import (
	"sync"
	"time"
)

// Allow at most 100 requests in the last 60 seconds
type SlidingWindowRateLimiter struct {
	mu       sync.Mutex
	limit    int
	requests []time.Time
	interval time.Duration
}

func NewSlidingWindowRateLimiter(limit int, interval time.Duration) RateLimiter {
	return &SlidingWindowRateLimiter{
		mu:       sync.Mutex{},
		limit:    limit,
		requests: make([]time.Time, 0, limit),
		interval: interval,
	}
}

func (rl *SlidingWindowRateLimiter) Allow() bool {
	// we clean up first for each round
	// then if length is still allows, add more requests
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// clean up
	now := time.Now()
	requests := make([]time.Time, 0, rl.limit)
	for _, log := range rl.requests {
		threshold := now.Add(-rl.interval)
		if log.After(threshold) {
			requests = append(requests, log)
		}
	}
	rl.requests = requests
	if len(rl.requests) < rl.limit {
		//allow
		requests = append(requests, now)
		rl.requests = requests
		return true
	}
	return false
}

// Cons
// Memory intensive: Must store all request timestamps -
// Slower: Must scan and clean old requests -
// Still allows all requests at once. With a 100 req/minute limit,
// you can receive all 100 requests in the first second, then wait 59 seconds.
//  It counts accurately but doesn't distribute evenly.
// Only enforces total count, not request distribution. Doesn't smooth traffic or enforce minimum intervals between requests.
