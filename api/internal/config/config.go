package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port    string
	AppName string
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

	return Config{
		Port:    port,
		AppName: appName,
	}
}
