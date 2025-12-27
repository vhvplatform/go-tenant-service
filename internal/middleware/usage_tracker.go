package middleware

import (
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// UsageTracker tracks API usage metrics
type UsageTracker struct {
	metrics map[string]*tenantMetrics
}

type tenantMetrics struct {
	apiCalls  int64
	bandwidth int64
}

// NewUsageTracker creates a new usage tracker
func NewUsageTracker() *UsageTracker {
	return &UsageTracker{
		metrics: make(map[string]*tenantMetrics),
	}
}

// Middleware returns a Gin middleware that tracks usage
func (ut *UsageTracker) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			c.Next()
			return
		}

		start := time.Now()

		// Process request
		c.Next()

		// Track metrics
		duration := time.Since(start)
		responseSize := c.Writer.Size()

		ut.recordMetrics(tenantID, 1, int64(responseSize))

		// Add metrics to response headers
		c.Header("X-Request-Duration", duration.String())
		c.Header("X-Response-Size", string(rune(responseSize)))
	}
}

func (ut *UsageTracker) recordMetrics(tenantID string, apiCalls int64, bandwidth int64) {
	if _, exists := ut.metrics[tenantID]; !exists {
		ut.metrics[tenantID] = &tenantMetrics{}
	}

	atomic.AddInt64(&ut.metrics[tenantID].apiCalls, apiCalls)
	atomic.AddInt64(&ut.metrics[tenantID].bandwidth, bandwidth)
}

// GetMetrics returns current metrics for a tenant
func (ut *UsageTracker) GetMetrics(tenantID string) (apiCalls int64, bandwidth int64) {
	if metrics, exists := ut.metrics[tenantID]; exists {
		return atomic.LoadInt64(&metrics.apiCalls), atomic.LoadInt64(&metrics.bandwidth)
	}
	return 0, 0
}
