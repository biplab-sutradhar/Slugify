package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RateLimitMiddleware implements token bucket rate limiting (100 requests/min per API key).
func RateLimitMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyID := c.GetString("api_key_id")
		if apiKeyID == "" {
			c.JSON(401, gin.H{"error": "API key required for rate limiting"})
			c.Abort()
			return
		}

		key := fmt.Sprintf("rate_limit:%s", apiKeyID)
		ctx := context.Background()

		// Get current tokens
		tokens, err := redisClient.Get(ctx, key).Float64()
		if err != nil {
			tokens = 100 // Reset on error
		}

		// Refill tokens (1 token per second, max 100)
		now := time.Now().Unix()
		lastRefill, _ := redisClient.Get(ctx, key+":last_refill").Int64()
		if now > lastRefill {
			refill := float64(now - lastRefill) // 1 token per second
			tokens = min(tokens+refill, 100)
			redisClient.Set(ctx, key, tokens, 0)
			redisClient.Set(ctx, key+":last_refill", now, 0)
		}

		if tokens < 1 {
			c.JSON(429, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		// Consume token
		tokens -= 1
		redisClient.Set(ctx, key, tokens, 0)

		c.Next()
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
