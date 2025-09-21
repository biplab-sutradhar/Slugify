package services

import (
	"context"
	"fmt"
	"github.com/biplab-sutradhar/slugify/api/internal/cache"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/idgen"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
	"time"
)

// LinkService handles business logic for links.
type LinkService struct {
	repo  db.LinkRepository
	cache cache.Cache
}

// NewLinkService creates a new service instance.
func NewLinkService(repo db.LinkRepository, cache cache.Cache) *LinkService {
	return &LinkService{repo: repo, cache: cache}
}

// SaveLink creates a new link with a generated short code and saves it to the repository.
func (s *LinkService) SaveLink(longURL string) (models.Link, error) {
	// Basic validation (non-empty URL)
	if longURL == "" {
		return models.Link{}, fmt.Errorf("long_url cannot be empty")
	}

	// Generate short code
	shortCode, err := idgen.GenerateShortCode()
	if err != nil {
		return models.Link{}, fmt.Errorf("failed to generate short code: %w", err)
	}
	// Create the link object
	link := models.Link{
		ID:        uuid.New().String(),
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	// Save the link to the database
	if err := s.repo.CreateLink(link); err != nil {
		return models.Link{}, fmt.Errorf("failed to save link: %w", err)
	}

	// Cache the URL (write-through)
	ctx := context.Background()
	if err := s.cache.SetURL(ctx, link.ShortCode, link.LongURL); err != nil {
		// Log cache error but continue (graceful degradation)
		fmt.Printf("Warning: Failed to cache URL: %v\n", err)
	}

	return link, nil
}

// GetLink retrieves a link by its short code.
// GetLink retrieves a link by its short code, checking cache first.
func (s *LinkService) GetLink(shortCode string) (models.Link, error) {
	ctx := context.Background()

	// Check cache
	longURL, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		// Log cache error but continue
		fmt.Printf("Warning: Cache error: %v\n", err)
	}
	if longURL != "" {
		// Cache hit: return a minimal Link struct
		return models.Link{
			ShortCode: shortCode,
			LongURL:   longURL,
		}, nil
	}

	// Cache miss: query database
	link, err := s.repo.GetLinkByShortCode(shortCode)
	if err != nil {
		return models.Link{}, err
	}

	// set Cache
	if err := s.cache.SetURL(ctx, link.ShortCode, link.LongURL); err != nil {
		fmt.Printf("Warning: Failed to cache URL: %v\n", err)
	}

	return link, nil
}
