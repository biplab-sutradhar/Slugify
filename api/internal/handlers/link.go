package handlers

import (
	"database/sql"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ShortenLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		link, err := service.SaveLink(req.LongURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp := models.ShortenResponse{
			ShortURL: "http://localhost:8080/" + link.ShortCode,
		}
		c.JSON(http.StatusCreated, resp)
	}
}

// ResolveLink handles GET /{shortCode} requests.
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

		c.Redirect(http.StatusFound, link.LongURL)
	}
}
