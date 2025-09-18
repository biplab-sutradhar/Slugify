package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
)

// Global DB variable (set in main.go)
var DB *sql.DB

// SaveLink inserts a new link into the database.
func SaveLink(longURL string) (models.Link, error) {
	// Generate a unique ID and short code
	link := models.Link{
		ID:        uuid.New().String(),
		ShortCode: fmt.Sprintf("%d", time.Now().UnixNano()), // temporary short code
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	// Save to DB
	if err := db.CreateLink(DB, link); err != nil {
		return models.Link{}, err
	}
	return link, nil
}

// GetLink fetches a link by short code from the database.
func GetLink(shortCode string) (models.Link, error) {
	return db.GetLinkByShortCode(DB, shortCode)
}
