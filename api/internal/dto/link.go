package dto

type ShortenRequest struct {
	LongURL string `json:"long_url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type UpdateLinkRequest struct {
	IsActive *bool `json:"is_active"`
}

type LinkResponse struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
	IsActive  bool   `json:"is_active"`
	Clicks    int64  `json:"clicks"`
	CreatedAt string `json:"created_at"` // or time.Time — match your API today
}
