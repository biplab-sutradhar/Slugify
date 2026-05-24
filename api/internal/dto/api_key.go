package dto

import "time"

// CreateAPIKeyRequest for POST /api/keys.
type CreateAPIKeyRequest struct {
	Name  string `json:"name"`
	Scope string `json:"scope"`
}

// APIKeyResponse for GET /api/keys.
type APIKeyResponse struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	Scope     string    `json:"scope"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	Usage     int64     `json:"usage"`
}
