//go:build ignore

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
)

func main() {
	connStr := "postgres://postgres:mynew@localhost:5432/urlshortener?sslmode=disable"
	database, err := db.NewDB(connStr)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	repo := db.NewPostgresLinkRepository(database)

	shortCode := "abc123"
	newLink := models.Link{
		ID:        uuid.New().String(),
		ShortCode: shortCode,
		LongURL:   "<https://example.com>",
		CreatedAt: time.Now(),
	}

	if err := repo.CreateLink(newLink); err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	fmt.Println("Inserted link:", newLink.ShortCode)

	fetched, err := repo.GetLinkByShortCode(shortCode)
	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}
	fmt.Println("Retrieved link:", fetched.LongURL)
}
