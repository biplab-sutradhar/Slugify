package handlers

import (
	"net/http"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
)

func CreateAPIKey(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
			return
		}

		var req models.CreateAPIKeyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if req.Name == "" {
			req.Name = "API key"
		}
		if req.Scope == "" {
			req.Scope = "default"
		}

		key, err := apiKeyService.CreateAPIKey(c, userID, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.APIKeyResponse{
			ID:        key.ID,
			Key:       key.Key,
			Name:      key.Name,
			Scope:     key.Scope,
			IsActive:  key.IsActive,
			CreatedAt: key.CreatedAt,
			Usage:     key.Usage,
		})
	}
}

func ListAPIKeys(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
			return
		}

		keys, err := apiKeyService.ListAPIKeys(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := make([]models.APIKeyResponse, 0, len(keys))
		for _, k := range keys {
			resp = append(resp, models.APIKeyResponse{
				ID:        k.ID,
				Key:       k.Key,
				Name:      k.Name,
				Scope:     k.Scope,
				IsActive:  k.IsActive,
				CreatedAt: k.CreatedAt,
				Usage:     k.Usage,
			})
		}
		c.JSON(http.StatusOK, resp)
	}
}

func DeleteAPIKey(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
			return
		}
		id := c.Param("id")
		if err := apiKeyService.DeleteAPIKey(c, id, userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
