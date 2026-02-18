package handler

import (
	"errors"
	"net/http"

	"github.com/dmitrycr/ShortUrl/internal/service"
	"github.com/go-chi/chi/v5"
)

// GetStats обрабатывает GET /api/stats/{code}
// Возвращает статистику по короткой ссылке
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		h.respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	stats, err := h.service.GetStats(r.Context(), shortCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrURLNotFound):
			h.respondError(w, http.StatusNotFound, "short URL not found")
		default:
			h.respondError(w, http.StatusInternalServerError, "failed to get stats")
		}
		return
	}

	h.respondJSON(w, http.StatusOK, stats)
}

// Delete обрабатывает DELETE /api/urls/{code}
// Удаляет короткую ссылку
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "code")
	if shortCode == "" {
		h.respondError(w, http.StatusBadRequest, "short code is required")
		return
	}

	err := h.service.DeleteURL(r.Context(), shortCode)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrURLNotFound):
			h.respondError(w, http.StatusNotFound, "short URL not found")
		default:
			h.respondError(w, http.StatusInternalServerError, "failed to delete URL")
		}
		return
	}

	h.respondJSON(w, http.StatusOK, SuccessResponse{
		Message: "URL deleted successfully",
	})
}
