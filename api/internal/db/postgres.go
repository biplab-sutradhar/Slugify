package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
)

// PostgresLinkRepository implements LinkRepository for PostgreSQL.
type PostgresLinkRepository struct {
	db *sql.DB
}

// NewPostgresLinkRepository creates a new repository instance.
func NewPostgresLinkRepository(db *sql.DB) *PostgresLinkRepository {
	return &PostgresLinkRepository{db: db}
}

// CreateLink inserts a link into the database.
func (r *PostgresLinkRepository) CreateLink(link models.Link) error {
	query := `
		INSERT INTO links (id, short_code, long_url, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, link.ID, link.ShortCode, link.LongURL, link.CreatedAt)
	return err
}

// GetLinkByShortCode retrieves a link by its short code.
func (r *PostgresLinkRepository) GetLinkByShortCode(shortCode string) (models.Link, error) {
	var link models.Link
	query := `
		SELECT id, short_code, long_url, created_at
		FROM links
		WHERE short_code = $1
	`
	err := r.db.QueryRow(query, shortCode).Scan(&link.ID, &link.ShortCode, &link.LongURL, &link.CreatedAt)
	if err == sql.ErrNoRows {
		return models.Link{}, err
	}
	if err != nil {
		return models.Link{}, err
	}
	return link, nil
}

// PostgresTicketRepository implements TicketRepository for PostgreSQL.
type PostgresTicketRepository struct {
	db *sql.DB
}

// NewPostgresTicketRepository creates a new ticket repository instance.
func NewPostgresTicketRepository(db *sql.DB) *PostgresTicketRepository {
	return &PostgresTicketRepository{db: db}
}

// SeedRanges inserts initial ranges into the database.
func (r *PostgresTicketRepository) SeedRanges(ctx context.Context, ranges []Range) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	for _, rng := range ranges {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO ranges (start_id, end_id, current_id, is_active)
			VALUES ($1, $2, $3, $4)
		`, rng.StartID, rng.EndID, rng.CurrentID, rng.IsActive)
		if err != nil {
			return fmt.Errorf("failed to seed range: %v", err)
		}
	}

	return tx.Commit()
}

// GetActiveRanges retrieves all active ranges.
func (r *PostgresTicketRepository) GetActiveRanges(ctx context.Context) ([]Range, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT range_id, start_id, end_id, current_id
		FROM ranges
		WHERE is_active = true
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query ranges: %v", err)
	}
	defer rows.Close()

	var ranges []Range
	for rows.Next() {
		var rng Range
		if err := rows.Scan(&rng.RangeID, &rng.StartID, &rng.EndID, &rng.CurrentID); err != nil {
			return nil, fmt.Errorf("failed to scan range: %v", err)
		}
		rng.IsActive = true
		ranges = append(ranges, rng)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading ranges: %v", err)
	}

	return ranges, nil
}

// LockAndIncrementRange locks a range and increments its current_id atomically.
func (r *PostgresTicketRepository) LockAndIncrementRange(ctx context.Context, rangeID int64) (int64, bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	var currentID, endID int64
	err = tx.QueryRowContext(ctx, `
		SELECT current_id, end_id
		FROM ranges
		WHERE range_id = $1 AND is_active = true
		FOR UPDATE
	`, rangeID).Scan(&currentID, &endID)
	if err == sql.ErrNoRows {
		return 0, false, fmt.Errorf("range %d is no longer active", rangeID)
	}
	if err != nil {
		return 0, false, fmt.Errorf("failed to lock range: %v", err)
	}

	if currentID >= endID {
		_, err = tx.ExecContext(ctx, `
			UPDATE ranges
			SET is_active = false
			WHERE range_id = $1
		`, rangeID)
		if err != nil {
			return 0, false, fmt.Errorf("failed to deactivate range: %v", err)
		}
		return 0, false, fmt.Errorf("range %d exhausted", rangeID)
	}

	newID := currentID + 1
	_, err = tx.ExecContext(ctx, `
		UPDATE ranges
		SET current_id = $1
		WHERE range_id = $2
	`, newID, rangeID)
	if err != nil {
		return 0, false, fmt.Errorf("failed to update current_id: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, false, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return newID, true, nil
}

// PostgresAPIKeyRepository implements APIKeyRepository for PostgreSQL.
type PostgresAPIKeyRepository struct {
	db *sql.DB
}

// NewPostgresAPIKeyRepository creates a new API key repository instance.
func NewPostgresAPIKeyRepository(db *sql.DB) *PostgresAPIKeyRepository {
	return &PostgresAPIKeyRepository{db: db}
}

// CreateAPIKey inserts a new API key into the database.
func (r *PostgresAPIKeyRepository) CreateAPIKey(ctx context.Context, key models.APIKey) error {
	query := `
		INSERT INTO api_keys (id, key, name, scope, is_active, created_at, usage)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(ctx, query, key.ID, key.Key, key.Name, key.Scope, key.IsActive, key.CreatedAt, key.Usage)
	return err
}

// GetAPIKeyByKey retrieves an API key by its key value.
func (r *PostgresAPIKeyRepository) GetAPIKeyByKey(ctx context.Context, key string) (models.APIKey, error) {
	var apiKey models.APIKey
	query := `
		SELECT id, key, name, scope, is_active, created_at, usage
		FROM api_keys
		WHERE key = $1
	`
	err := r.db.QueryRowContext(ctx, query, key).Scan(
		&apiKey.ID, &apiKey.Key, &apiKey.Name, &apiKey.Scope, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.Usage,
	)
	if err == sql.ErrNoRows {
		return models.APIKey{}, err
	}
	if err != nil {
		return models.APIKey{}, err
	}
	return apiKey, nil
}

// GetAPIKeys retrieves all API keys.
func (r *PostgresAPIKeyRepository) GetAPIKeys(ctx context.Context) ([]models.APIKey, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, key, name, scope, is_active, created_at, usage
		FROM api_keys
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []models.APIKey
	for rows.Next() {
		var key models.APIKey
		if err := rows.Scan(
			&key.ID, &key.Key, &key.Name, &key.Scope, &key.IsActive, &key.CreatedAt, &key.Usage,
		); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// DeleteAPIKey deletes an API key by ID.
func (r *PostgresAPIKeyRepository) DeleteAPIKey(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM api_keys WHERE id = $1", id)
	return err
}

// IncrementUsage increments the usage count for an API key.
func (r *PostgresAPIKeyRepository) IncrementUsage(ctx context.Context, apiKeyID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE api_keys
		SET usage = usage + 1
		WHERE id = $1
	`, apiKeyID)
	return err
}
