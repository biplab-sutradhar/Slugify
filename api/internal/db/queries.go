package db

import (
	"database/sql"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
	// "time"
)

// PostgresLinkRepository implements LinkRepository for PostgreSQL.
type PostgresLinkRepository struct {
	db *sql.DB
}

// NewPostgresLinkRepository creates a new repository instance.
func NewPostgresLinkRepository(db *sql.DB) *PostgresLinkRepository {
	return &PostgresLinkRepository{db: db}
}

func (r *PostgresLinkRepository) CreateLink(link models.Link) error {
	query := `INSERT INTO links (id, short_code, long_url, created_at)VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, link.ID, link.ShortCode, link.LongURL, link.CreatedAt)
	return err
}

func (r *PostgresLinkRepository) GetLinkByShortCode(shortCode string) (models.Link, error) {
	var link models.Link
	query := `
		SELECT id, short_code, long_url, created_at
		FROM links
		WHERE short_code = $1
	`
	err := r.db.QueryRow(query, shortCode).Scan(&link.ID, &link.ShortCode, &link.LongURL, &link.CreatedAt)

	if err != nil {
		return models.Link{}, err
	}
	return link, nil
}
