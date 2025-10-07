package main

import (
	"log"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/cache"
	"github.com/biplab-sutradhar/slugify/api/internal/config"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/handlers"
	"github.com/biplab-sutradhar/slugify/api/internal/idgen"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	database, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Println("Connected to Postgres")
	}
	defer database.Close()

	// Initialize Redis connection
	redisClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}
	defer redisClient.Close()

	// Initialize Ticket server connection
	ticketServer, err := idgen.NewTicketServer(database)
	if err != nil {
		log.Fatalf("Failed to connect to ticketServer: %v", err)
	}

	// Initialize repository and service
	repo := db.NewPostgresLinkRepository(database)
	service := services.NewLinkService(repo, redisClient, ticketServer)

	// Set up Gin router
	r := gin.Default()

	// Middleware to log request duration
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
	})

	// Register endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/api/shorten", handlers.ShortenLink(service))
	r.GET("/:shortCode", handlers.ResolveLink(service))

	// Start server
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
