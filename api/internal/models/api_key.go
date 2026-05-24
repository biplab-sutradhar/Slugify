package models

import "time"

// APIKey represents an API key in the database.
type APIKey struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Key       string    `json:"key" db:"key"`
	Name      string    `json:"name" db:"name"`
	Scope     string    `json:"scope" db:"scope"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Usage     int64     `json:"usage" db:"usage"`
}
