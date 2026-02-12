package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrycr/ShortUrl/internal/model"
	"github.com/dmitrycr/ShortUrl/internal/storage"
	"github.com/dmitrycr/ShortUrl/internal/validator"
	"github.com/dmitrycr/ShortUrl/pkg/generator"
)

var (
	ErrURLNotFound     = errors.New("url not found")
	ErrURLExpired      = errors.New("url has expired")
	ErrInvalidURL      = errors.New("invalid url")
	ErrCodeAlreadyUsed = errors.New("short code already in use")
)

type URLService struct {
	storage   storage.Storage
	generator *generator.Generator
	validator *validator.URLValidator
	baseURL   string
}

type Config struct {
	Storage    storage.Storage
	BaseURL    string
	CodeLength int
}

func NewURLService(cfg Config) *URLService {
	codeLength := cfg.CodeLength
	if codeLength == 0 {
		codeLength = 6
	}

	return &URLService{
		storage:   cfg.Storage,
		generator: generator.NewGenerator(codeLength),
		validator: validator.NewURLValidator(),
		baseURL:   cfg.BaseURL,
	}
}

func (s *URLService) ShortenURL(ctx context.Context, req *model.CreateURLRequest) (*model.CreateURLResponse, error) {
	normalizedURL := s.validator.NormalizeURL(req.URL)

	if err := s.validator.ValidateURL(req.URL); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	var (
		shortCode string
		err       error
	)

	if req.CustomCode != "" {
		if err := s.validator.ValidateCustomCode(req.CustomCode); err != nil {
			return nil, fmt.Errorf("invalid custom code: %w", err)
		}
		shortCode = req.CustomCode
	} else {
		shortCode, err = s.generateUniqueCode(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate code: %w", err)
		}
	}
	// Вычисляем время истечения
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(req.ExpiresIn) * time.Second)
		expiresAt = &expiry
	}

	url := &model.URL{
		OriginalURL: normalizedURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now(),
		ExpiresAt:   expiresAt,
		ClickCount:  0,
	}

	// save on storage
	err = s.storage.Save(ctx, url)
	if err != nil {
		if errors.Is(err, storage.ErrDuplicateCode) {
			return nil, ErrCodeAlreadyUsed
		}
		return nil, fmt.Errorf("failed to save url: %w", err)
	}

	return &model.CreateURLResponse{
		ShortURL:    s.buildShortURL(shortCode),
		ShortCode:   shortCode,
		OriginalURL: normalizedURL,
		ExpiresAt:   expiresAt,
	}, nil
}

func (s *URLService) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.storage.GetByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return "", ErrURLNotFound
		}
		if errors.Is(err, storage.ErrExpired) {
			return "", ErrURLExpired
		}
		return "", fmt.Errorf("failed to get url: %w", err)
	}

	return url.OriginalURL, nil
}

func (s *URLService) RegisterClick(ctx context.Context, shortCode string) error {
	_, err := s.storage.GetByShortCode(ctx, shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return ErrURLNotFound
		}
		if errors.Is(err, storage.ErrExpired) {
			return ErrURLExpired
		}
		return fmt.Errorf("failed to get url: %w", err)
	}

	//увеличиваем счетчик
	if err := s.storage.IncrementClicks(ctx, shortCode); err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	return nil
}

func (s *URLService) GetStats(ctx context.Context, shortCode string) (*model.Stats, error) {
	stats, err := s.storage.GetStats(ctx, shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return nil, ErrURLNotFound
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	return stats, nil
}

func (s *URLService) DeleteURL(ctx context.Context, shortCode string) error {
	if err := s.storage.Delete(ctx, shortCode); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to delete url: %w", err)
	}
	return nil
}

// generateUniqueCode генерирует уникальный короткий код
func (s *URLService) generateUniqueCode(ctx context.Context) (string, error) {
	const maxAttempts = 5

	for attempt := 0; attempt < maxAttempts; attempt++ {
		code, err := s.generator.Generate()
		if err != nil {
			return "", err
		}

		// Проверяем уникальность
		_, err = s.storage.GetByShortCode(ctx, code)
		if errors.Is(err, storage.ErrNotFound) {
			// Код свободен!
			return code, nil
		}

		// Код занят, пробуем еще раз
	}

	return "", errors.New("failed to generate unique code after multiple attempts")
}

func (s *URLService) buildShortURL(code string) string {
	return s.baseURL + "/" + code
}
