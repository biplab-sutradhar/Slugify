package handlers

import (
	"net/http"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateAPIKey handles POST /api/keys or POST /auth/keys.
func CreateAPIKey(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Support both API key auth and JWT auth
		userID := c.GetString("user_id")
		if userID == "" {
			// If no JWT user_id, get from API key auth context if available
			userID = c.GetString("api_key_user_id")
		}

		var req models.CreateAPIKeyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		key, err := apiKeyService.CreateAPIKey(c, userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := models.APIKeyResponse{
			ID:        key.ID,
			Key:       key.Key,
			Name:      key.Name,
			Scope:     key.Scope,
			IsActive:  key.IsActive,
			CreatedAt: key.CreatedAt,
			Usage:     key.Usage,
		}
		c.JSON(http.StatusCreated, resp)
	}
}

// ListAPIKeys handles GET /api/keys or GET /auth/keys.
func ListAPIKeys(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Support both API key auth and JWT auth
		userID := c.GetString("user_id")
		if userID == "" {
			// If no JWT user_id, get from API key auth context if available
			userID = c.GetString("api_key_user_id")
		}

		keys, err := apiKeyService.ListAPIKeys(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var resp []models.APIKeyResponse
		for _, key := range keys {
			resp = append(resp, models.APIKeyResponse{
				ID:        key.ID,
				Key:       key.Key,
				Name:      key.Name,
				Scope:     key.Scope,
				IsActive:  key.IsActive,
				CreatedAt: key.CreatedAt,
				Usage:     key.Usage,
			})
		}
		c.JSON(http.StatusOK, resp)
	}
}

// DeleteAPIKey handles DELETE /api/keys/{id}.
func DeleteAPIKey(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := apiKeyService.DeleteAPIKey(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
