// Package api holds HTTP handlers: decode the request, call a usecase, encode
// the response. It contains ZERO business logic.
package api

import (
	"net/http"

	"gestion-commande/internal/service"
)

// Handler adapts HTTP requests to the service usecases.
type Handler struct {
	svc *service.Service
}

// NewHandler wires the handler to its usecase dependency.
func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// ListenAndServe registers routes and starts the HTTP server.
func (h *Handler) ListenAndServe(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.health)
	return http.ListenAndServe(addr, mux)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
