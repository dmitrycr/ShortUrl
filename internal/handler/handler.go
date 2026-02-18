package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dmitrycr/ShortUrl/internal/service"
)

type Handler struct {
	service *service.URLService
	logger  *slog.Logger
}

type ErrorResponse struct {
	Error  string `json:"error"`
	Status int    `json:"status"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func New(service *service.URLService, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, ErrorResponse{
		Error:  message,
		Status: status,
	})
}

func (h *Handler) decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(dst)
}
