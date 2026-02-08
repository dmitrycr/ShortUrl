package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dmitrycr/ShortUrl/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, connString string) (*PostgresStorage, error) {
	// Настройка пула для подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	//Настройка пула
	config.MaxConns = 25                      // Максимум подключений
	config.MinConns = 5                       // Минимум подключений
	config.MaxConnLifetime = time.Hour        // Время жизни подключения
	config.MaxConnIdleTime = 30 * time.Minute // Время простоя
	config.HealthCheckPeriod = time.Minute    // Период проверки здоровья

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresStorage{
		pool: pool,
	}, nil
}

func (s *PostgresStorage) Save(ctx context.Context, url *model.URL) error {
	query := `
		INSERT INFO urls (original_url, short_code, created_at, expires_at, click_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := s.pool.QueryRow(
		ctx,
		query,
		url.OriginalURL,
		url.ShortCode,
		url.CreatedAt,
		url.ExpiresAt,
		url.ClickCount,
	).Scan(&url.ID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateCode
		}
		return fmt.Errorf("failed to save url: %w", err)
	}

	return nil
}

func (s *PostgresStorage) GetByShortCode(ctx context.Context, code string) (*model.URL, error) {
	query := `
        SELECT id, original_url, short_code, created_at, expires_at, click_count
        FROM urls
        WHERE short_code = $1
    `
	var url model.URL

	err := s.pool.QueryRow(ctx, query, code).Scan(
		&url.ID,
		&url.OriginalURL,
		&url.ShortCode,
		&url.CreatedAt,
		&url.ExpiresAt,
		&url.ClickCount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get url: %w", err)
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, ErrExpired
	}

	return &url, nil
}

func (s *PostgresStorage) IncrementClicks(ctx context.Context, code string) error {
	query := `
	UPDATE urls
	SET click_count = click_count + 1
	WHERE short_code = $1
`
	result, err := s.pool.Exec(ctx, query, code)
	if err != nil {
		return fmt.Errorf("failed to increment clicks: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStorage) GetStats(ctx context.Context, code string) (*model.Stats, error) {
	query := `
		SELECT short_code, original_url, click_count, created_at, expires_at
		FROM urls
		WHERE short_code = $1
	`
	var stats model.Stats

	err := s.pool.QueryRow(ctx, query, code).Scan(
		&stats.ShortCode,
		&stats.OriginalURL,
		&stats.ClickCount,
		&stats.CreatedAt,
		&stats.ExpiresAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &stats, nil

}

func (s *PostgresStorage) Delete(ctx context.Context, code string) error {
	query := `
		DELETE FROM urls 
		WHERE short_code = $1
`

	result, err := s.pool.Exec(ctx, query, code)
	if err != nil {
		return fmt.Errorf("failed delete url: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PostgresStorage) Close() error {
	s.pool.Close()
	return nil
}
