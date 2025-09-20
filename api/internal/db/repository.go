package db

import "github.com/biplab-sutradhar/slugify/api/internal/models"

type LinkRepository interface {
	CreateLink(link models.Link) error
	GetLinkByShortCode(ShortCode string) (models.Link, error)
}
