package dto

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type MintAPIKeyRequest struct {
	Name string `json:"name"`
}
