package storage

import (
	"context"
	"errors"

	"github.com/dmitrycr/ShortUrl/internal/model"
)

var (
	ErrNotFound      = errors.New("url not found")
	ErrDuplicateCode = errors.New("short code already exists")
	ErrExpired       = errors.New("url has expired")
)

type Storage interface {
	Save(ctx context.Context, url *model.URL) error
	GetByShortCode(ctx context.Context, code string) (*model.URL, error)
	IncrementClicks(ctx context.Context, code string) error
	GetStats(ctx context.Context, code string) (*model.Stats, error)
	Delete(ctx context.Context, code string) error
	Close() error
}
