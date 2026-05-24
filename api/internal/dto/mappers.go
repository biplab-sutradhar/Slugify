package dto

import (
	"github.com/biplab-sutradhar/slugify/api/internal/models"
)
import "time"

func UserFromModel(u models.User) UserResponse {
	return UserResponse{
		ID: u.ID, Email: u.Email, Name: u.Name, CreatedAt: u.CreatedAt,
	}
}

func LinkFromModel(l models.Link) LinkResponse {
	return LinkResponse{
		ID: l.ID, ShortCode: l.ShortCode, LongURL: l.LongURL,
		IsActive: l.IsActive, Clicks: l.Clicks, CreatedAt: l.CreatedAt.Format(time.RFC3339),
	}
}

func APIKeyFromModel(k models.APIKey) APIKeyResponse {
	return APIKeyResponse{
		ID: k.ID, Key: k.Key, Name: k.Name, Scope: k.Scope,
		IsActive: k.IsActive, CreatedAt: k.CreatedAt, Usage: k.Usage,
	}
}
