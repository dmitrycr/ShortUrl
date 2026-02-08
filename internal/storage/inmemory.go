package storage

import (
	"context"
	"sync"
	"time"

	"github.com/dmitrycr/ShortUrl/internal/model"
)

// InMemoryStorage реализует Storage в памяти для тестов
type InMemoryStorage struct {
	mu     sync.RWMutex
	urls   map[string]*model.URL // short_code -> URL
	nextID int64
}

// NewInMemoryStorage создает новое in-memory хранилище
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls:   make(map[string]*model.URL),
		nextID: 1,
	}
}

// Save сохраняет URL в памяти
func (s *InMemoryStorage) Save(ctx context.Context, url *model.URL) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем дубликат
	if _, exists := s.urls[url.ShortCode]; exists {
		return ErrDuplicateCode
	}

	// Присваиваем ID
	url.ID = s.nextID
	s.nextID++

	// Копируем для защиты от изменений
	urlCopy := *url
	s.urls[url.ShortCode] = &urlCopy

	return nil
}

// GetByShortCode получает URL по коду
func (s *InMemoryStorage) GetByShortCode(ctx context.Context, code string) (*model.URL, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urls[code]
	if !exists {
		return nil, ErrNotFound
	}

	// Проверяем истечение
	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpired
	}

	// Копируем для защиты от изменений
	urlCopy := *url
	return &urlCopy, nil
}

// IncrementClicks увеличивает счетчик
func (s *InMemoryStorage) IncrementClicks(ctx context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	url, exists := s.urls[code]
	if !exists {
		return ErrNotFound
	}

	url.ClickCount++
	return nil
}

// GetStats возвращает статистику
func (s *InMemoryStorage) GetStats(ctx context.Context, code string) (*model.Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urls[code]
	if !exists {
		return nil, ErrNotFound
	}

	return &model.Stats{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		ClickCount:  url.ClickCount,
		CreatedAt:   url.CreatedAt,
		ExpiresAt:   url.ExpiresAt,
	}, nil
}

// Delete удаляет URL
func (s *InMemoryStorage) Delete(ctx context.Context, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.urls[code]; !exists {
		return ErrNotFound
	}

	delete(s.urls, code)
	return nil
}

// Close ничего не делает для in-memory
func (s *InMemoryStorage) Close() error {
	return nil
}

// Clear очищает все данные (для тестов)
func (s *InMemoryStorage) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.urls = make(map[string]*model.URL)
	s.nextID = 1
}
