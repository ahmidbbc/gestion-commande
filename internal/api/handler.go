// Package api holds HTTP handlers: decode the request, call a usecase, encode
// the response. It contains ZERO business logic.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gestion-commande/internal/service"
	"gestion-commande/internal/store"
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
	mux.HandleFunc("GET /commandes", h.listCommandes)
	mux.HandleFunc("GET /commandes/{id}", h.getCommande)
	mux.HandleFunc("PATCH /commandes/{id}/statut", h.updateStatut)
	return http.ListenAndServe(addr, mux)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

type updateStatutRequest struct {
	Statut service.Statut `json:"statut"`
}

type commandeResponse struct {
	ID     string         `json:"id"`
	Statut service.Statut `json:"statut"`
	CreeLe string         `json:"cree_le"`
}

func (h *Handler) listCommandes(w http.ResponseWriter, r *http.Request) {
	statut := service.Statut(r.URL.Query().Get("statut"))
	commandes, err := h.svc.ListCommandes(statut)
	switch {
	case errors.Is(err, service.ErrStatutInvalide):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case err != nil:
		http.Error(w, "erreur interne", http.StatusInternalServerError)
		return
	}

	out := make([]commandeResponse, 0, len(commandes))
	for _, c := range commandes {
		out = append(out, commandeResponse{
			ID:     c.ID,
			Statut: c.Statut,
			CreeLe: c.CreeLe.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *Handler) getCommande(w http.ResponseWriter, r *http.Request) {
	c, err := h.svc.GetCommande(r.PathValue("id"))
	switch {
	case errors.Is(err, store.ErrCommandeIntrouvable):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "erreur interne", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(commandeResponse{
		ID:     c.ID,
		Statut: c.Statut,
		CreeLe: c.CreeLe.Format(time.RFC3339),
	})
}

func (h *Handler) updateStatut(w http.ResponseWriter, r *http.Request) {
	var req updateStatutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corps de requête invalide", http.StatusBadRequest)
		return
	}

	c, err := h.svc.UpdateStatut(r.PathValue("id"), req.Statut)
	switch {
	case errors.Is(err, service.ErrStatutInvalide):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.Is(err, store.ErrCommandeIntrouvable):
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, "erreur interne", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(commandeResponse{
		ID:     c.ID,
		Statut: c.Statut,
		CreeLe: c.CreeLe.Format(time.RFC3339),
	})
}
