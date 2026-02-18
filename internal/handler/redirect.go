package handler

import (
	"errors"
	"net/http"

	"github.com/dmitrycr/ShortUrl/internal/service"
	"github.com/go-chi/chi/v5"
)

// Redirect обрабатывает GET /{code}
// Перенаправляет пользователя на оригинальный URL
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Получаем код из URL
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		h.respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	// Получаем оригинальный URL
	originalURL, err := h.service.GetOriginalURL(r.Context(), shortCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrURLNotFound):
			h.respondError(w, http.StatusNotFound, "short URL not found")
		case errors.Is(err, service.ErrURLExpired):
			h.respondError(w, http.StatusGone, "this short URL has expired")
		default:
			h.respondError(w, http.StatusInternalServerError, "failed to resolve URL")
		}
		return
	}

	// Регистрируем клик асинхронно — не задерживаем редирект
	go func() {
		if err := h.service.RegisterClick(r.Context(), shortCode); err != nil {
			h.logger.Error("failed to register click",
				"code", shortCode,
				"error", err,
			)
		}
	}()

	// 301 — постоянный редирект (кешируется браузером)
	// 302 — временный редирект (не кешируется)
	// Для счетчика кликов лучше 302
	http.Redirect(w, r, originalURL, http.StatusFound)
}
