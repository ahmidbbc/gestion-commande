package service

import (
	"errors"
	"testing"
)

// fakeRepo is a stub Repository (health port).
type fakeRepo struct{}

func (fakeRepo) Ping() error { return nil }

// fakeCommandeRepo is a configurable in-test CommandeRepository.
type fakeCommandeRepo struct {
	all      []Commande
	byStatut []Commande
	errList  error
}

func (f *fakeCommandeRepo) Create(Commande) error        { return nil }
func (f *fakeCommandeRepo) Save(Commande) error          { return nil }
func (f *fakeCommandeRepo) Get(string) (Commande, error) { return Commande{}, nil }
func (f *fakeCommandeRepo) List() ([]Commande, error) {
	return f.all, f.errList
}
func (f *fakeCommandeRepo) ListByStatut(Statut) ([]Commande, error) {
	return f.byStatut, f.errList
}

func TestListCommandes_Succes(t *testing.T) {
	repo := &fakeCommandeRepo{
		all: []Commande{
			{ID: "1", Statut: StatutEnCours},
			{ID: "2", Statut: StatutLivree},
		},
	}
	svc := New(fakeRepo{}, repo)

	got, err := svc.ListCommandes("")
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("attendu 2 commandes, obtenu %d", len(got))
	}
}

func TestListCommandes_StatutInvalide(t *testing.T) {
	svc := New(fakeRepo{}, &fakeCommandeRepo{})

	_, err := svc.ListCommandes("inconnu")
	if !errors.Is(err, ErrStatutInvalide) {
		t.Fatalf("attendu ErrStatutInvalide, obtenu %v", err)
	}
}

func TestListCommandes_ErreurStore(t *testing.T) {
	boom := errors.New("store indisponible")
	svc := New(fakeRepo{}, &fakeCommandeRepo{errList: boom})

	_, err := svc.ListCommandes("")
	if !errors.Is(err, boom) {
		t.Fatalf("attendu l'erreur du store, obtenu %v", err)
	}
}
