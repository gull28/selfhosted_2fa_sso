package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	rateLimiters    map[string]*rate.Limiter
	limiterInterval time.Duration
	limiterRequests int
	mu              sync.Mutex
}

func NewRateLimiter(interval time.Duration, requests int) *RateLimiter {
	return &RateLimiter{
		rateLimiters:    make(map[string]*rate.Limiter),
		limiterInterval: interval,
		limiterRequests: requests,
	}
}

func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if limiter, exists := rl.rateLimiters[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Every(rl.limiterInterval/time.Duration(rl.limiterRequests)), rl.limiterRequests)
	rl.rateLimiters[ip] = limiter
	return limiter
}
