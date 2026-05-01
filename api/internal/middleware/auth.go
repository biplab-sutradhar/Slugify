package middleware

import (
	"net/http"

	"github.com/biplab-sutradhar/slugify/api/internal/auth"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware authenticates requests using API keys
// and exposes the owning user via the Gin context.
func AuthMiddleware(apiKeyRepo db.APIKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			return
		}

		if err := auth.ValidateAPIKey(apiKey); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key format"})
			return
		}

		key, err := apiKeyRepo.GetAPIKeyByKey(c, apiKey)
		if err != nil || !key.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or inactive API key"})
			return
		}

		c.Set("api_key_id", key.ID)
		c.Set("user_id", key.UserID)
		c.Next()
	}
}
