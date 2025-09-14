package main

import (
	"fmt"
	"github.com/biplab-sutradhar/slugify/api/internal/config"
	"github.com/biplab-sutradhar/slugify/api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	router.POST("/api/shorten", handlers.ShortenLink)
	router.GET("/:shortCode", handlers.ResolveLink)
	addr := ":" + cfg.Port
	fmt.Printf("Starting server on %s...\n", addr)
	router.Run(addr)
}
