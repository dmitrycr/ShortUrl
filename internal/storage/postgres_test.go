package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/dmitrycr/ShortUrl/internal/model"
)

// TestPostgresStorage_Integration запускается только при наличии TEST_DATABASE_URL
func TestPostgresStorage_Integration(t *testing.T) {
	connString := os.Getenv("TEST_DATABASE_URL")
	if connString == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration tests")
	}

	ctx := context.Background()

	// Создаем подключение
	storage, err := NewPostgresStorage(ctx, connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer storage.Close()

	t.Run("Save and Get", func(t *testing.T) {
		url := &model.URL{
			OriginalURL: "https://example.com/test",
			ShortCode:   "test123",
			CreatedAt:   time.Now(),
			ClickCount:  0,
		}

		// Сохраняем
		err := storage.Save(ctx, url)
		if err != nil {
			t.Fatalf("Save failed: %v", err)
		}

		// Получаем
		retrieved, err := storage.GetByShortCode(ctx, "test123")
		if err != nil {
			t.Fatalf("GetByShortCode failed: %v", err)
		}

		if retrieved.OriginalURL != url.OriginalURL {
			t.Errorf("Expected %s, got %s", url.OriginalURL, retrieved.OriginalURL)
		}

		// Очищаем
		storage.Delete(ctx, "test123")
	})

	t.Run("Duplicate ShortCode", func(t *testing.T) {
		url := &model.URL{
			OriginalURL: "https://example.com/dup",
			ShortCode:   "dup123",
			CreatedAt:   time.Now(),
		}

		// Первое сохранение
		err := storage.Save(ctx, url)
		if err != nil {
			t.Fatalf("First save failed: %v", err)
		}

		// Второе сохранение с тем же кодом
		url2 := &model.URL{
			OriginalURL: "https://example.com/another",
			ShortCode:   "dup123",
			CreatedAt:   time.Now(),
		}

		err = storage.Save(ctx, url2)
		if err != ErrDuplicateCode {
			t.Errorf("Expected ErrDuplicateCode, got %v", err)
		}

		// Очищаем
		storage.Delete(ctx, "dup123")
	})

	t.Run("IncrementClicks", func(t *testing.T) {
		url := &model.URL{
			OriginalURL: "https://example.com/clicks",
			ShortCode:   "clicks123",
			CreatedAt:   time.Now(),
			ClickCount:  0,
		}

		storage.Save(ctx, url)

		// Увеличиваем счетчик
		for i := 0; i < 5; i++ {
			err := storage.IncrementClicks(ctx, "clicks123")
			if err != nil {
				t.Fatalf("IncrementClicks failed: %v", err)
			}
		}

		// Проверяем
		stats, err := storage.GetStats(ctx, "clicks123")
		if err != nil {
			t.Fatalf("GetStats failed: %v", err)
		}

		if stats.ClickCount != 5 {
			t.Errorf("Expected 5 clicks, got %d", stats.ClickCount)
		}

		storage.Delete(ctx, "clicks123")
	})

	t.Run("Expired URL", func(t *testing.T) {
		expiredTime := time.Now().Add(-1 * time.Hour)
		url := &model.URL{
			OriginalURL: "https://example.com/expired",
			ShortCode:   "expired123",
			CreatedAt:   time.Now(),
			ExpiresAt:   &expiredTime,
		}

		storage.Save(ctx, url)

		// Пытаемся получить истекшую ссылку
		_, err := storage.GetByShortCode(ctx, "expired123")
		if err != ErrExpired {
			t.Errorf("Expected ErrExpired, got %v", err)
		}

		storage.Delete(ctx, "expired123")
	})
}
