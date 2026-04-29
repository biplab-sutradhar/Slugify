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

type LinkService struct {
	repo         db.LinkRepository
	cache        cache.Cache
	ticketServer idgen.TicketServer
	apiKeyRepo   db.APIKeyRepository
	domainURL    string
}

func NewLinkService(repo db.LinkRepository, cache cache.Cache, ticketServer idgen.TicketServer, apiKeyRepo db.APIKeyRepository, domainURL string) *LinkService {
	return &LinkService{repo: repo, cache: cache, ticketServer: ticketServer, apiKeyRepo: apiKeyRepo, domainURL: domainURL}
}

func (s *LinkService) GetDomainURL() string {
	return s.domainURL
}

func (s *LinkService) SaveLink(longURL string) (models.Link, error) {
	if longURL == "" {
		return models.Link{}, fmt.Errorf("long_url cannot be empty")
	}

	ctx := context.Background()
	shortCode, err := s.ticketServer.GenerateID(ctx)
	if err != nil {
		return models.Link{}, fmt.Errorf("failed to generate short code: %v", err)
	}

	link := models.Link{
		ID:        uuid.New().String(),
		ShortCode: shortCode,
		LongURL:   longURL,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateLink(link); err != nil {
		return models.Link{}, err
	}

	if err := s.cache.SetURL(ctx, link.ShortCode, link.LongURL); err != nil {
		fmt.Printf("Warning: Failed to cache URL: %v\\n", err)
	}

	return link, nil
}

func (s *LinkService) IncrementAPIKeyUsage(ctx context.Context, apiKeyID string) error {
	return s.apiKeyRepo.IncrementUsage(ctx, apiKeyID)
}

func (s *LinkService) GetLink(shortCode string) (models.Link, error) {
	ctx := context.Background()

	longURL, err := s.cache.GetURL(ctx, shortCode)
	if err != nil {
		fmt.Printf("Warning: Cache error: %v\\n", err)
	}
	if longURL != "" {
		return models.Link{
			ShortCode: shortCode,
			LongURL:   longURL,
			IsActive:  true,
		}, nil
	}

	link, err := s.repo.GetLinkByShortCode(shortCode)
	if err != nil {
		return models.Link{}, err
	}

	if !link.IsActive {
		return models.Link{}, fmt.Errorf("link is deactivated")
	}

	if err := s.cache.SetURL(ctx, link.ShortCode, link.LongURL); err != nil {
		fmt.Printf("Warning: Failed to cache URL: %v\\n", err)
	}

	return link, nil
}

func (s *LinkService) GetLinkByID(id string) (models.Link, error) {
	return s.repo.GetLinkByID(id)
}

func (s *LinkService) ListLinks(limit, offset int) ([]models.Link, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.ListLinks(limit, offset)
}

func (s *LinkService) UpdateLinkStatus(id string, isActive bool) error {
	return s.repo.UpdateLinkStatus(id, isActive)
}

func (s *LinkService) DeleteLink(id string) error {
	return s.repo.DeleteLink(id)
}
func (s *LinkService) IncrementClicks(shortCode string) {
	if err := s.repo.IncrementClicks(shortCode); err != nil {
		fmt.Printf("Warning: Failed to increment clicks: %v\\n", err)
	}
}
