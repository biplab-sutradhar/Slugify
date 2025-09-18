package handlers

import (
	"database/sql"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ShortenLink handles the request to shorten a given long URL.
func ShortenLink(c *gin.Context) {
	var req models.ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil || !strings.HasPrefix(req.LongURL, "http") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	link, err := services.SaveLink(req.LongURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save link"})
		return
	}

	c.JSON(http.StatusCreated, models.ShortenResponse{
		ShortURL: "http://localhost:9000/" + link.ShortCode,
	})
}

// ResolveLink handles resolving a short URL to its long URL.
func ResolveLink(c *gin.Context) {
	shortCode := c.Param("shortCode")

	link, err := services.GetLink(shortCode)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.Redirect(http.StatusFound, link.LongURL)
}
