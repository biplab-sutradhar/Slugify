package handlers

import (
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// ShortenLink handles the request to shorten a given long URL.
func ShortenLink(c *gin.Context) {
	var req models.ShortenRequest

	// Bind the incoming JSON body to the ShortenRequest struct.
	if err := c.ShouldBindJSON(&req); err != nil || !strings.HasPrefix(req.LongURL, "http") { // Ensure the URL starts with "http".
		// If binding fails or URL is invalid, return an error response.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing URL"})
		return
	}

	// Save the long URL and generate a shortened link.
	link := services.SaveLink(req.LongURL)

	// Return the short URL in the response.
	c.JSON(http.StatusCreated, models.ShortenResponse{
		ShortURL: "http://localhost:9000/" + link.ShortCode,
	})
}

// ResolveLink handles the request to resolve a short URL to its original long URL.
func ResolveLink(c *gin.Context) {
	// Retrieve the short code from the URL parameter.
	shortCode := c.Param("shortCode")

	// Lookup the long URL associated with the short code.
	link, exists := services.GetLink(shortCode)
	if !exists {
		// If the short code does not exist, return a 404 error.
		c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
		return
	}

	// Redirect the user to the original long URL.
	c.Redirect(http.StatusFound, link.LongURL)
}
