// Package api holds HTTP handlers: decode the request, call a usecase, encode
// the response. It contains ZERO business logic.
package api

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"gestion-commande/internal/service"
)

//go:embed templates/home.html
var templatesFS embed.FS

var homeTmpl = template.Must(template.New("home.html").Funcs(template.FuncMap{
	"statusClass": statusClass,
}).ParseFS(templatesFS, "templates/home.html"))

// statusClass maps an order status to its CSS badge modifier.
func statusClass(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "payée":
		return "badge-payée"
	case "en cours":
		return "badge-encours"
	case "expédiée":
		return "badge-expédiée"
	case "annulée":
		return "badge-annulée"
	default:
		return "badge-default"
	}
}

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
	mux.HandleFunc("/", h.home)
	mux.HandleFunc("/health", h.health)
	return http.ListenAndServe(addr, mux)
}

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	orders, err := h.svc.ListOrders()
	if err != nil {
		http.Error(w, "impossible de charger les commandes", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := homeTmpl.Execute(w, struct{ Orders []service.Order }{Orders: orders}); err != nil {
		http.Error(w, "erreur de rendu", http.StatusInternalServerError)
	}
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
