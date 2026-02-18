package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter создает и настраивает роутер
func NewRouter(h *Handler) http.Handler {
	r := chi.NewRouter()

	// -- Middleware --

	// Логирование каждого запроса
	r.Use(middleware.Logger)

	// Восстановление после паники
	r.Use(middleware.Recoverer)

	// Таймаут на запрос
	r.Use(middleware.Timeout(30 * time.Second))

	// Сжатие ответов
	r.Use(middleware.Compress(5))

	// Заголовки безопасности
	r.Use(securityHeaders)

	// -- Роуты --

	// Health check
	r.Get("/health", h.Health)

	// API группа
	r.Route("/api", func(r chi.Router) {
		// Создание короткой ссылки
		r.Post("/shorten", h.Shorten)

		// Статистика
		r.Get("/stats/{code}", h.GetStats)

		// Удаление
		r.Delete("/urls/{code}", h.Delete)
	})

	// Редирект — должен быть последним
	r.Get("/{code}", h.Redirect)

	return r
}

// securityHeaders добавляет базовые заголовки безопасности
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
