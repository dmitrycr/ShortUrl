package model

import "time"

type URL struct {
	ID          int64      `db:"id"`
	OriginalURL string     `db:"original_url"`
	ShortCode   string     `db:"short_url"`
	CreatedAt   time.Time  `db:"created_at"`
	ExpiresAt   *time.Time `db:"expires_at"`
	ClickCount  int64      `db:"click_count"`
}

// Stats - статистика оп ссылке
type Stats struct {
	OriginalURL string     `json:"original_url"`
	ShortCode   string     `json:"short_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	ClickCount  int64      `json:"click_count"`
}

// CreateURLRequest - запрос на создание короткой ссылки
type CreateURLRequest struct {
	URL        string `json:"url" validate:"required,url"`
	CustomCode string `json:"custom_code,omitempty"` // опционально
	ExpiresIn  int    `json:"expires_in,omitempty"`  // В секундах
}

// CreateURLResponse - ответ при создании короткой ссылки
type CreateURLResponse struct {
	ShortURL    string     `json:"short_url"`
	ShortCode   string     `json:"short_code"`
	OriginalURL string     `json:"original_url"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}
