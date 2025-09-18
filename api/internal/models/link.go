package models

import "time"

type Link struct {
	ID        string    `json:"id" db:"id"`
	ShortCode string    `json:"short_code" db:"short_code"`
	LongURL   string    `json:"long_url" db:"long_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}
