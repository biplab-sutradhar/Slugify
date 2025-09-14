package handlers

import (
	"net/http"
	"strings"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
)

func ShortenLink(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.ShouldBindJSON(&req); err != nil || !strings.HasPrefix(req.LongURL, "http") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	link := services.SaveLink(req.LongURL)

	c.JSON(http.StatusCreated, models.ShortenResponse{
		ShortURL: "http://localhost:9000/" + link.ShortCode,
	})
}

func ResolveLink(c *gin.Context) {
	shortCode := c.Param("shortCode")
	link, exists := services.GetLink(shortCode)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
		return
	}

	c.Redirect(http.StatusFound, link.LongURL)
}
