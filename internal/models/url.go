package models

import "time"

type URLRequest struct {
	LongURL string `json:"long_url"`
}

type URL struct {
	ID        uint64    `json:"id"`
	ShortURL  string    `json:"short_url"`
	LongURL   string    `json:"long_url"`
	UserID    uint64    `json:"user_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetAllURLsResponse struct {
	URLs []*URLResponse `json:"urls"`
}

type URLResponse struct {
	ID       uint64 `json:"id"`
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
	Clicks   uint64 `json:"clicks"`
}
