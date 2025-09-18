package main

import (
	"log"

	"github.com/biplab-sutradhar/slugify/api/internal/config"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/handlers"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// connect DB
	database, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	// set global DB
	services.DB = database

	// setup router
	r := gin.Default()
	r.POST("/api/shorten", handlers.ShortenLink)
	r.GET("/:shortCode", handlers.ResolveLink)

	log.Printf("Server running on :%s", cfg.Port)
	r.Run(":" + cfg.Port)
}
