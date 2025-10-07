package services

import (
	"context"
	"fmt"
	"time"

	"github.com/biplab-sutradhar/slugify/api/internal/cache"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/idgen"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
)

// LinkService handles business logic for links.
type LinkService struct {
	repo         db.LinkRepository
	cache        cache.Cache
	ticketServer idgen.TicketServer
}

// NewLinkService creates a new service instance.
func NewLinkService(repo db.LinkRepository, cache cache.Cache, ticketServer idgen.TicketServer) *LinkService {
	return &LinkService{repo: repo, cache: cache, ticketServer: ticketServer}
}

// SaveLink creates a new link with a generated short code and caches it.
func (s *LinkService) SaveLink(longURL string) (models.Link, error) {
	// Basic validation
	if longURL == "" {
		return models.Link{}, fmt.Errorf("long_url cannot be empty")
	}

	// Generate short code
	ctx := context.Background()
	shortCode, err := s.ticketServer.GenerateID(ctx)
	if err != nil {
		return models.Link{}, fmt.Errorf("failed to generate short code: %v", err)
	}

	// Create link
	link := models.Link{
		ID:        uuid.New().String(),
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := s.repo.CreateLink(link); err != nil {
		return models.Link{}, err
	}

	// Cache the URL (write-through)
	if err := s.cache.SetURL(ctx, link.ShortCode, link.LongURL); err != nil {
		// Log cache error but continue
		fmt.Printf("Warning: Failed to cache URL: %v\n", err)
	}

	return link, nil
}

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

	// Cache the result
	if err := s.cache.SetURL(ctx, link.ShortCode, link.LongURL); err != nil {
		// Log cache error but continue
		fmt.Printf("Warning: Failed to cache URL: %v\n", err)
	}

	return link, nil
}
