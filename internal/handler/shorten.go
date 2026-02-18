package handler

import (
	"errors"
	"net/http"

	"github.com/dmitrycr/ShortUrl/internal/model"
	"github.com/dmitrycr/ShortUrl/internal/service"
)

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	// Декодируем тело запроса
	var req model.CreateURLRequest
	if err := h.decodeJSON(r, &req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Передаем в сервис
	resp, err := h.service.ShortenURL(r.Context(), &req)
	if err != nil {
		h.logger.Error("failed to shorten url",
			"url", req.URL,
			"error", err,
		)

		switch {
		case errors.Is(err, service.ErrInvalidURL):
			h.respondError(w, http.StatusBadRequest, "invalid URL provided")
		case errors.Is(err, service.ErrCodeAlreadyUsed):
			h.respondError(w, http.StatusConflict, "this custom code is already taken")
		default:
			h.respondError(w, http.StatusInternalServerError, "failed to shorten URL")
		}
		return
	}

	h.respondJSON(w, http.StatusCreated, resp)
}
