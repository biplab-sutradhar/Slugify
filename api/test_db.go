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

	// Create sample link
	newLink := models.Link{
		ID:        uuid.New().String(),
		ShortCode: "abc12",
		LongURL:   "https://example.com",
		CreatedAt: time.Now(),
	}

	// Insert into DB
	if err := db.CreateLink(database, newLink); err != nil {
		log.Fatalf("Insert failed: %v", err)
	}
	fmt.Println("✅ Inserted link:", newLink.ShortCode)

	// Fetch from DB
	fetched, err := db.GetLinkByShortCode(database, "abc123")
	if err != nil {
		log.Fatalf("Fetch failed: %v", err)
	}
	fmt.Println("🔎 Retrieved link:", fetched.LongURL)
}
