package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple in-memory rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per duration
	duration time.Duration
}

// Visitor tracks request count and last reset time
type Visitor struct {
	count    int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, duration time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		duration: duration,
	}

	// Clean up old visitors periodically
	go rl.cleanup()

	return rl
}

// cleanup removes old visitor entries
func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(rl.duration)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.duration*2 {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &Visitor{
			count:    1,
			lastSeen: time.Now(),
		}
		return true
	}

	// Reset if duration has passed
	if time.Since(v.lastSeen) > rl.duration {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	// Check if rate limit exceeded
	if v.count >= rl.rate {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

// Global rate limiter (100 requests per minute)
var globalLimiter = NewRateLimiter(100, time.Minute)

// RateLimitMiddleware applies rate limiting
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !globalLimiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

