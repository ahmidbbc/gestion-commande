package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gestion-commande/internal/service"
)

// fakeRepo satisfies service.Repository (health port).
type fakeRepo struct{}

func (fakeRepo) Ping() error { return nil }

// fakeCommandeRepo is a configurable service.CommandeRepository for tests.
type fakeCommandeRepo struct {
	all       []service.Commande
	byStatut  []service.Commande
	errList   error
	errCreate error
}

func (f *fakeCommandeRepo) Create(service.Commande) error        { return f.errCreate }
func (f *fakeCommandeRepo) Save(service.Commande) error          { return nil }
func (f *fakeCommandeRepo) Get(string) (service.Commande, error) { return service.Commande{}, nil }
func (f *fakeCommandeRepo) List() ([]service.Commande, error)    { return f.all, f.errList }
func (f *fakeCommandeRepo) ListByStatut(service.Statut) ([]service.Commande, error) {
	return f.byStatut, f.errList
}

func newTestHandler(commandes service.CommandeRepository) *Handler {
	return NewHandler(service.New(fakeRepo{}, commandes))
}

func TestCreateCommande_Succes(t *testing.T) {
	repo := &fakeCommandeRepo{}
	body := `{"client":"Alice","produits":[{"nom":"Café","quantite":2,"prix_unite":150}]}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/commandes", strings.NewReader(body))

	newTestHandler(repo).createCommande(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("attendu 201, obtenu %d", rec.Code)
	}
	var out commandeResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("réponse JSON invalide: %v", err)
	}
	if out.Client != "Alice" {
		t.Fatalf("attendu client Alice, obtenu %q", out.Client)
	}
	if out.Montant != 300 {
		t.Fatalf("attendu montant 300, obtenu %d", out.Montant)
	}
	if out.Statut != service.StatutEnCours {
		t.Fatalf("attendu statut en_cours, obtenu %q", out.Statut)
	}
	if out.ID == "" {
		t.Fatal("attendu un ID généré, obtenu vide")
	}
}

func TestCreateCommande_CorpsInvalide(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/commandes", strings.NewReader("{"))

	newTestHandler(&fakeCommandeRepo{}).createCommande(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("attendu 400, obtenu %d", rec.Code)
	}
}

func TestCreateCommande_ValidationEchoue(t *testing.T) {
	body := `{"client":"","produits":[{"nom":"Café","quantite":1,"prix_unite":150}]}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/commandes", strings.NewReader(body))

	newTestHandler(&fakeCommandeRepo{}).createCommande(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("attendu 400, obtenu %d", rec.Code)
	}
}

func TestCreateCommande_ErreurStore(t *testing.T) {
	repo := &fakeCommandeRepo{errCreate: errors.New("store indisponible")}
	body := `{"client":"Alice","produits":[{"nom":"Café","quantite":1,"prix_unite":150}]}`
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/commandes", strings.NewReader(body))

	newTestHandler(repo).createCommande(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("attendu 500, obtenu %d", rec.Code)
	}
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
