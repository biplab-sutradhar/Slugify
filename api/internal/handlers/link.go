package handlers

import (
	"database/sql"

	"net/http"
	"strconv"

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

		go service.IncrementClicks(shortCode)

		c.Redirect(http.StatusFound, link.LongURL)
	}
}

func GetLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		link, err := service.GetLinkByID(id)
		if err == sql.ErrNoRows {
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

		links, err := service.ListLinks(limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, links)
	}
}

func UpdateLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if err := service.UpdateLinkStatus(id, *req.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Link updated"})
	}
}

func DeleteLink(service *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeleteLink(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
