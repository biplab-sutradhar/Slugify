package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port        string
	AppName     string
	DatabaseURL string
	RedisURL    string
	DomainURL   string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set in environment variables")
	}

	appName := os.Getenv("APP_NAME")
	if appName == "" {
		appName = "url-shortener"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		log.Fatal("REDIS_URL is not set in environment variables")
	}

	domainURL := os.Getenv("DOMAINURL")
	if domainURL == "" {
		log.Fatal("DOMAINURL is not set in environment variables")
	}

	return Config{
		Port:        port,
		AppName:     appName,
		DatabaseURL: dbURL,
		RedisURL:    redisUrl,
		DomainURL:   domainURL,
	}
}
