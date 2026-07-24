package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gestion-commande/internal/service"
)

// fakeRepo satisfies service.Repository (health port).
type fakeRepo struct{}

func (fakeRepo) Ping() error { return nil }

// fakeCommandeRepo is a configurable service.CommandeRepository for tests.
type fakeCommandeRepo struct {
	all      []service.Commande
	byStatut []service.Commande
	errList  error
}

func (f *fakeCommandeRepo) Create(service.Commande) error        { return nil }
func (f *fakeCommandeRepo) Save(service.Commande) error          { return nil }
func (f *fakeCommandeRepo) Get(string) (service.Commande, error) { return service.Commande{}, nil }
func (f *fakeCommandeRepo) List() ([]service.Commande, error)    { return f.all, f.errList }
func (f *fakeCommandeRepo) ListByStatut(service.Statut) ([]service.Commande, error) {
	return f.byStatut, f.errList
}

func newTestHandler(commandes service.CommandeRepository) *Handler {
	return NewHandler(service.New(fakeRepo{}, commandes))
}

func TestListCommandes_Succes(t *testing.T) {
	repo := &fakeCommandeRepo{
		all: []service.Commande{
			{ID: "1", Statut: service.StatutEnCours},
			{ID: "2", Statut: service.StatutLivree},
		},
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/commandes", nil)

	newTestHandler(repo).listCommandes(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("attendu 200, obtenu %d", rec.Code)
	}
	var out []commandeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("réponse JSON invalide: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("attendu 2 commandes, obtenu %d", len(out))
	}
}

func TestListCommandes_StatutInvalide(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/commandes?statut=inconnu", nil)

	newTestHandler(&fakeCommandeRepo{}).listCommandes(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("attendu 400, obtenu %d", rec.Code)
	}
}

func TestListCommandes_ErreurStore(t *testing.T) {
	repo := &fakeCommandeRepo{errList: errors.New("store indisponible")}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/commandes", nil)

	newTestHandler(repo).listCommandes(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("attendu 500, obtenu %d", rec.Code)
	}
}
