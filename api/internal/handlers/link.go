package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
)

func ShortenLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		apiKeyID := c.GetString("api_key_id")
		if apiKeyID != "" {
			if err := service.IncrementAPIKeyUsage(c, apiKeyID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update usage"})
				return
			}
		}

		link, err := service.SaveLink(req.LongURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := models.ShortenResponse{
			ShortURL: service.GetDomainURL() + "/" + link.ShortCode,
		}
		c.JSON(http.StatusCreated, resp)
	}
}

func ResolveLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		link, err := service.GetLink(shortCode)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		apiKeyID := c.GetString("api_key_id")
		if apiKeyID != "" {
			if err := service.IncrementAPIKeyUsage(c, apiKeyID); err != nil {
				fmt.Printf("Warning: Failed to update usage: %v\\n", err)
			}
		}

		c.Redirect(http.StatusFound, link.LongURL)
	}
}
