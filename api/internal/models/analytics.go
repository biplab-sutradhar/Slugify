package models

import "time"

type ClickEvent struct {
	ID        int64     `json:"id" db:"id"`
	LinkID    string    `json:"link_id" db:"link_id"`
	ShortCode string    `json:"short_code" db:"short_code"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	IPHash    string    `json:"ip_hash" db:"ip_hash"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referrer  string    `json:"referrer" db:"referrer"`
	Country   string    `json:"country" db:"country"`
}

type LinkAnalytics struct {
	TotalClicks  int64           `json:"total_clicks"`
	ClicksByDay  []DayClickCount `json:"clicks_by_day"`
	TopReferrers []ReferrerCount `json:"top_referrers"`
	TopCountries []CountryCount  `json:"top_countries"`
}

type DayClickCount struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type ReferrerCount struct {
	Referrer string `json:"referrer"`
	Count    int64  `json:"count"`
}

type CountryCount struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}
