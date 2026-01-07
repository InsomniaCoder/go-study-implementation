package main

import (
	"sync"
	"time"
)

type RateLimiter interface {
	Allow() bool
}

type FixedWindowLimiter struct {
	limit    int
	interval time.Duration
	mu       sync.Mutex
	count    int
	window   time.Time
}

func NewFixedWindowLimiter(interval time.Duration, limit int) RateLimiter {
	return &FixedWindowLimiter{
		limit:    limit,
		mu:       sync.Mutex{},
		count:    0,
		interval: interval,
		window:   time.Now(), //start time
	}
}

func (rl *FixedWindowLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if rl.window.Add(rl.interval).After(now) {
		//within the window
		if rl.count < rl.limit {
			rl.count++
			return true
		} else {
			return false
		}
	} else {
		rl.count = 1
		rl.window = now
		return true
	}
}

//Cons
// 	Burst problem at boundaries: 100 requests at 11:59:59 + 100 at 12:00:01 = 200 requests in 2 seconds (double the rate)
//  Unfair: Users timing requests at window edges can exceed intended limits
//  Less accurate: Rate enforcement is coarse-grained
