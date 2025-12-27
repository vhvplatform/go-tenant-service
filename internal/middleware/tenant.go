package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// TenantContext extracts tenant information from headers and adds to context
func TenantContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		tenantTier := c.GetHeader("X-Tenant-Tier")

		if tenantID != "" {
			c.Set("tenant_id", tenantID)
		}

		if tenantTier != "" {
			c.Set("tenant_tier", tenantTier)
		}

		c.Next()
	}
}

// APIKeyAuth validates API keys for tenant authentication
func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			// Try Authorization header with Bearer token
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		// Skip auth for health endpoints
		if strings.HasPrefix(c.Request.URL.Path, "/health") || strings.HasPrefix(c.Request.URL.Path, "/ready") {
			c.Next()
			return
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "API key required",
			})
			c.Abort()
			return
		}

		// Here you would validate the API key against your database
		// For now, we'll just set it in context
		c.Set("api_key", apiKey)
		c.Next()
	}
}

// TenantIsolation ensures tenant data isolation
func TenantIsolation() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")

		// Validate tenant ID is present for non-public endpoints
		if tenantID == "" && !isPublicEndpoint(c.Request.URL.Path) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Tenant identification required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func isPublicEndpoint(path string) bool {
	publicPaths := []string{"/health", "/ready", "/api/v1/tenants"}
	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	return false
}
