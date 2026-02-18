package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmitrycr/ShortUrl/internal/config"
	"github.com/dmitrycr/ShortUrl/internal/handler"
	"github.com/dmitrycr/ShortUrl/internal/service"
	"github.com/dmitrycr/ShortUrl/internal/storage"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Настраиваем логгер
	logger := setupLogger(cfg)
	logger.Info("starting url shortener service",
		"environment", cfg.Environment,
		"port", cfg.ServerPort,
	)

	// Подключаемся к базе данных
	ctx := context.Background()
	store, err := storage.NewPostgresStorage(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	logger.Info("connected to database")

	// Создаем сервис
	urlService := service.NewURLService(service.Config{
		Storage:    store,
		BaseURL:    cfg.BaseURL,
		CodeLength: cfg.CodeLength,
	})

	// Создаем handlers
	h := handler.New(urlService, logger)

	// Создаем роутер
	router := handler.NewRouter(h)

	// Создаем HTTP сервер
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		logger.Info("server is listening", "address", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Таймаут для завершения
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}

// setupLogger настраивает логгер в зависимости от окружения
func setupLogger(cfg *config.Config) *slog.Logger {
	var handler slog.Handler

	if cfg.IsDevelopment() {
		// В dev режиме — красивый текстовый вывод
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		// В production — JSON формат
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	return slog.New(handler)
}
