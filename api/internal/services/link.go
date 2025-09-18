package services

import (
	"fmt"
	"github.com/biplab-sutradhar/slugify/api/internal/models"
	"sync"
	"time"
)

var (
	links   sync.Map // thread-safe map
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
	links.Store(shortCode, link)
	return link
}

func GetLink(shortCode string) (models.Link, bool) {
	value, exists := links.Load(shortCode)
	if exists {
		return value.(models.Link), true
	}
	return models.Link{}, false
}
