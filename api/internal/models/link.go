package models

import "time"

type Link struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ShortCode string    `json:"short_code" db:"short_code"`
	LongURL   string    `json:"long_url" db:"long_url"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Clicks    int64     `json:"clicks" db:"clicks"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type UpdateLinkRequest struct {
	IsActive *bool `json:"is_active"`
}
