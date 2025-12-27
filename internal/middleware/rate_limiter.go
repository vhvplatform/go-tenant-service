package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	limits map[string]*tierLimit
	mu     sync.RWMutex
}

type tierLimit struct {
	requestsPerMinute int
	tokens            int
	lastRefill        time.Time
	mu                sync.Mutex
}

// NewRateLimiter creates a new rate limiter with tier-based limits
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limits: map[string]*tierLimit{
			"free":         {requestsPerMinute: 60, tokens: 60, lastRefill: time.Now()},
			"basic":        {requestsPerMinute: 300, tokens: 300, lastRefill: time.Now()},
			"professional": {requestsPerMinute: 1000, tokens: 1000, lastRefill: time.Now()},
			"enterprise":   {requestsPerMinute: 5000, tokens: 5000, lastRefill: time.Now()},
		},
	}
}

// Middleware returns a Gin middleware function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tier := c.GetString("tenant_tier")
		if tier == "" {
			tier = "free"
		}

		if !rl.allowRequest(tier) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded for your subscription tier",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) allowRequest(tier string) bool {
	rl.mu.RLock()
	limit, exists := rl.limits[tier]
	rl.mu.RUnlock()

	if !exists {
		return false
	}

	limit.mu.Lock()
	defer limit.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(limit.lastRefill)
	tokensToAdd := int(elapsed.Minutes() * float64(limit.requestsPerMinute))

	if tokensToAdd > 0 {
		limit.tokens = min(limit.requestsPerMinute, limit.tokens+tokensToAdd)
		limit.lastRefill = now
	}

	if limit.tokens > 0 {
		limit.tokens--
		return true
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
