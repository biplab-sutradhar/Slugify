package db

import (
	"context"

	"github.com/biplab-sutradhar/slugify/api/internal/models"
)

type LinkRepository interface {
	CreateLink(link models.Link) error
	GetLinkByShortCode(shortCode string) (models.Link, error)
}

// Range represents an ID range in the database.
type Range struct {
	RangeID   int64
	StartID   int64
	EndID     int64
	CurrentID int64
	IsActive  bool
}

// TicketRepository defines methods for managing ID ranges.
type TicketRepository interface {
	SeedRanges(ctx context.Context, ranges []Range) error
	GetActiveRanges(ctx context.Context) ([]Range, error)
	LockAndIncrementRange(ctx context.Context, rangeID int64) (int64, bool, error)
}

// APIKeyRepository defines methods for API key management.
type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, key models.APIKey) error
	GetAPIKeyByKey(ctx context.Context, key string) (models.APIKey, error)
	GetAPIKeys(ctx context.Context) ([]models.APIKey, error)
	DeleteAPIKey(ctx context.Context, id string) error
	IncrementUsage(ctx context.Context, apiKeyID string) error
}
