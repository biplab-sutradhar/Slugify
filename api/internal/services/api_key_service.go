package services

import (
	"context"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/auth"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
)

// APIKeyService handles API key business logic.
type APIKeyService struct {
	repo db.APIKeyRepository
}

// NewAPIKeyService creates a new service instance.
func NewAPIKeyService(repo db.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

// CreateAPIKey creates a new API key.
func (s *APIKeyService) CreateAPIKey(ctx context.Context, req models.CreateAPIKeyRequest) (models.APIKey, error) {
	key, err := auth.GenerateAPIKey()
	if err != nil {
		return models.APIKey{}, err
	}

	apiKey := models.APIKey{
		ID:        uuid.New().String(),
		Key:       key,
		Name:      req.Name,
		Scope:     req.Scope,
		IsActive:  true,
		CreatedAt: time.Now(),
		Usage:     0,
	}

	if err := s.repo.CreateAPIKey(ctx, apiKey); err != nil {
		return models.APIKey{}, err
	}

	return apiKey, nil
}

// ListAPIKeys lists all API keys.
func (s *APIKeyService) ListAPIKeys(ctx context.Context) ([]models.APIKey, error) {
	return s.repo.GetAPIKeys(ctx)
}

// DeleteAPIKey deletes an API key by ID.
func (s *APIKeyService) DeleteAPIKey(ctx context.Context, id string) error {
	return s.repo.DeleteAPIKey(ctx, id)
}
