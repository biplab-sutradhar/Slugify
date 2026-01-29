package middleware

import (
	"net/http"

	"github.com/biplab-sutradhar/slugify/api/internal/auth"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware authenticates requests using API keys.
func AuthMiddleware(apiKeyRepo db.APIKeyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		// Validate API key format
		if err := auth.ValidateAPIKey(apiKey); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key format"})
			c.Abort()
			return
		}

		// Check if API key exists and is active
		key, err := apiKeyRepo.GetAPIKeyByKey(c, apiKey)
		if err != nil || !key.IsActive {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or inactive API key"})
			c.Abort()
			return
		}

		// Set API key ID in context
		c.Set("api_key_id", key.ID)
		c.Next()
	}
}
