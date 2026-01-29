package idgen

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

// TicketServer defines the interface for ID generation.
type TicketServer interface {
	GenerateID(ctx context.Context) (string, error)
	Close() error
}

// PostgresTicketServer implements TicketServer using PostgreSQL.
type PostgresTicketServer struct {
	db *sql.DB
}

// Range represents an ID range in the database.
type Range struct {
	RangeID   int64
	StartID   int64
	EndID     int64
	CurrentID int64
	IsActive  bool
}

// NewTicketServer initializes the ticket server and seeds ranges if none exist.
func NewTicketServer(db *sql.DB) (*PostgresTicketServer, error) {
	// Check if ranges exist
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM ranges").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to check ranges: %v", err)
	}

	// Seed 4 ranges if none exist (10M IDs each)
	if count == 0 {
		_, err := db.Exec(`
			INSERT INTO ranges (start_id, end_id, current_id, is_active)
			VALUES
				(1000000, 11000000, 1000000, true),
				(11000001, 21000000, 11000001, true),
				(21000001, 31000000, 21000001, true),
				(31000001, 41000000, 31000001, true)
		`)
		if err != nil {
			return nil, fmt.Errorf("failed to seed ranges: %v", err)
		}
	}

	return &PostgresTicketServer{db: db}, nil
}

// Close is a no-op for now, as the database connection is managed externally.
func (ts *PostgresTicketServer) Close() error {
	return nil
}

// GenerateID generates a unique base62-encoded ID using a random range.
func (ts *PostgresTicketServer) GenerateID(ctx context.Context) (string, error) {
	// Get active ranges
	rows, err := ts.db.QueryContext(ctx, "SELECT range_id, start_id, end_id, current_id FROM ranges WHERE is_active = true")
	if err != nil {
		return "", fmt.Errorf("failed to query ranges: %v", err)
	}
	defer rows.Close()

	var ranges []Range
	for rows.Next() {
		var r Range
		if err := rows.Scan(&r.RangeID, &r.StartID, &r.EndID, &r.CurrentID); err != nil {
			return "", fmt.Errorf("failed to scan range: %v", err)
		}
		ranges = append(ranges, r)
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error reading ranges: %v", err)
	}
	if len(ranges) == 0 {
		return "", fmt.Errorf("no active ranges available")
	}

	// Select a random range
	rand.Seed(time.Now().UnixNano())
	selectedRange := ranges[rand.Intn(len(ranges))]

	// Start a transaction
	tx, err := ts.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Lock the selected range and increment current_id
	var newID int64
	err = tx.QueryRowContext(ctx, `
		SELECT current_id
		FROM ranges
		WHERE range_id = $1 AND is_active = true
		FOR UPDATE
	`, selectedRange.RangeID).Scan(&newID)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("range %d is no longer active", selectedRange.RangeID)
	}
	if err != nil {
		return "", fmt.Errorf("failed to lock range: %v", err)
	}

	if newID >= selectedRange.EndID {
		// Mark range as inactive
		_, err = tx.ExecContext(ctx, `
			UPDATE ranges
			SET is_active = false
			WHERE range_id = $1
		`, selectedRange.RangeID)
		if err != nil {
			return "", fmt.Errorf("failed to deactivate range: %v", err)
		}
		return "", fmt.Errorf("range %d exhausted", selectedRange.RangeID)
	}

	newID++
	// FIXED: Changed $1 to $2 for the WHERE clause
	_, err = tx.ExecContext(ctx, `
		UPDATE ranges
		SET current_id = $1
		WHERE range_id = $2
	`, newID, selectedRange.RangeID)
	if err != nil {
		return "", fmt.Errorf("failed to update current_id: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Encode ID to base62 with secret mapping
	encodedID := Encode(newID)
	return encodedID, nil
}
