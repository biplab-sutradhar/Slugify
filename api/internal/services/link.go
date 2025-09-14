package services

import (
	"fmt"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"time"
)

var (
	links   = make(map[string]models.Link)
	counter = 0
)

func SaveLink(longURL string) models.Link {
	counter++
	shortCode := fmt.Sprintf("%d", counter)
	link := models.Link{
		ShortCode: shortCode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}
	links[shortCode] = link
	return link
}

func GetLink(shortCode string) (models.Link, bool) {
	link, exists := links[shortCode]
	return link, exists
}
