package handlers

import (
	"net/http"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"

	"github.com/gin-gonic/gin"
)

// CreateAPIKey handles POST /api/keys.
func CreateAPIKey(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateAPIKeyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		key, err := apiKeyService.CreateAPIKey(c, req)
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

// ListAPIKeys handles GET /api/keys.
func ListAPIKeys(apiKeyService *services.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		keys, err := apiKeyService.ListAPIKeys(c)
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
