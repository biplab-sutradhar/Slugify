package services

import (
	"fmt"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
	"github.com/biplab-sutradhar/slugify/api/internal/idgen"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"github.com/google/uuid"
	"time"
)

// LinkService handles business logic for links.
type LinkService struct {
	repo db.LinkRepository
}

// NewLinkService creates a new service instance.
func NewLinkService(repo db.LinkRepository) *LinkService {
	return &LinkService{repo: repo}
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

	return link, nil
}

// GetLink retrieves a link by its short code.
func (s *LinkService) GetLink(shortCode string) (models.Link, error) {
	link, err := s.repo.GetLinkByShortCode(shortCode)
	if err != nil {
		return models.Link{}, fmt.Errorf("link not found: %w", err)
	}
	return link, nil
}
