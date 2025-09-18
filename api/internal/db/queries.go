package db

import (
	"database/sql"
	"errors"
	// "time"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
)

func CreateLink(db *sql.DB, link models.Link) error {
	query := `
        INSERT INTO links (id, short_code, long_url, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := db.Exec(query, link.ID, link.ShortCode, link.LongURL, link.CreatedAt)
	return err
}

func GetLinkByShortCode(db *sql.DB, shortCode string) (models.Link, error) {
	query := `
        SELECT id, short_code, long_url, created_at
        FROM links
        WHERE short_code = $1
    `
	var link models.Link
	err := db.QueryRow(query, shortCode).Scan(&link.ID, &link.ShortCode, &link.LongURL, &link.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return link, err
	}
	return link, err
}
