package main

import (
	"log"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/cache"
	"github.com/biplab-sutradhar/slugify/api/internal/config"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/handlers"
	"github.com/biplab-sutradhar/slugify/api/internal/idgen"
	"github.com/biplab-sutradhar/slugify/api/internal/middleware"
	"github.com/biplab-sutradhar/slugify/api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := config.LoadConfig()

	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migrations applied successfully")

	database, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Println("Connected to Postgres")
	}
	defer database.Close()

	redisClient, err := cache.NewRedisClient(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}
	defer redisClient.Close()

	ticketServer, err := idgen.NewTicketServer(database)
	if err != nil {
		log.Fatalf("Failed to connect to ticketServer: %v", err)
	} else {
		log.Println("Connected to ticketServer")
	}

	repo := db.NewPostgresLinkRepository(database)
	apiKeyRepo := db.NewPostgresAPIKeyRepository(database)
	linkService := services.NewLinkService(repo, redisClient, ticketServer, apiKeyRepo, cfg.DomainURL)
	apiKeyService := services.NewAPIKeyService(apiKeyRepo)
	userRepo := db.NewPostgresUserRepository(database)
	authService := services.NewAuthService(userRepo, apiKeyRepo, cfg.JWTSecret)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Request %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public redirect route (no auth needed)
	r.GET("/:shortCode", handlers.ResolveLink(linkService))

	// Public API routes to create API keys (no auth needed)
	// publicAPI := r.Group("/api")
	// {
	// 	publicAPI.POST("/keys", handlers.CreateAPIKey(apiKeyService))
	// }

	// Protected API routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(apiKeyRepo))
	api.Use(middleware.RateLimitMiddleware(redisClient.GetRawClient()))
	{
		api.POST("/shorten", handlers.ShortenLink(linkService))
		api.GET("/links", handlers.ListLinks(linkService))
		api.GET("/links/:id", handlers.GetLink(linkService))
		api.PATCH("/links/:id", handlers.UpdateLink(linkService))
		api.DELETE("/links/:id", handlers.DeleteLink(linkService))
		api.GET("/keys", handlers.ListAPIKeys(apiKeyService))
		api.POST("/keys", handlers.CreateAPIKey(apiKeyService))
		api.DELETE("/keys/:id", handlers.DeleteAPIKey(apiKeyService))
	}

	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register(authService))
		auth.POST("/login", handlers.Login(authService))
	}

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
