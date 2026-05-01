package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
)

func ShortenLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
			return
		}

		var req models.ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if apiKeyID := c.GetString("api_key_id"); apiKeyID != "" {
			if err := service.IncrementAPIKeyUsage(c, apiKeyID); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update usage"})
				return
			}
		}

		link, err := service.SaveLink(userID, req.LongURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, models.ShortenResponse{
			ShortURL: service.GetDomainURL() + "/" + link.ShortCode,
		})
	}
}

func ResolveLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		link, err := service.GetLink(shortCode)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short code not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		go service.IncrementClicks(shortCode)
		c.Redirect(http.StatusFound, link.LongURL)
	}
}

func GetLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		link, err := service.GetLinkByID(id, userID)
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, link)
	}
}

func ListLinks(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not identified"})
			return
		}

		limit := 20
		offset := 0
		if l := c.Query("limit"); l != "" {
			if parsed, err := strconv.Atoi(l); err == nil {
				limit = parsed
			}
		}
		if o := c.Query("offset"); o != "" {
			if parsed, err := strconv.Atoi(o); err == nil {
				offset = parsed
			}
		}

		links, err := service.ListLinks(userID, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if links == nil {
			links = []models.Link{}
		}
		c.JSON(http.StatusOK, links)
	}
}

func UpdateLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		var req models.UpdateLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		if req.IsActive == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "is_active is required"})
			return
		}
		if err := service.UpdateLinkStatus(id, userID, *req.IsActive); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Link updated"})
	}
}

func DeleteLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		if err := service.DeleteLink(id, userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
