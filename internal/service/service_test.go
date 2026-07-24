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
	all       []Commande
	byStatut  []Commande
	errList   error
	errCreate error
	created   *Commande
}

func (f *fakeCommandeRepo) Create(c Commande) error {
	if f.errCreate != nil {
		return f.errCreate
	}
	f.created = &c
	return nil
}
func (f *fakeCommandeRepo) Save(Commande) error          { return nil }
func (f *fakeCommandeRepo) Get(string) (Commande, error) { return Commande{}, nil }
func (f *fakeCommandeRepo) List() ([]Commande, error) {
	return f.all, f.errList
}
func (f *fakeCommandeRepo) ListByStatut(Statut) ([]Commande, error) {
	return f.byStatut, f.errList
}

func TestCreateCommande_Succes(t *testing.T) {
	repo := &fakeCommandeRepo{}
	svc := New(fakeRepo{}, repo)

	c, err := svc.CreateCommande("Alice", []Produit{
		{Nom: "Café", Quantite: 2, PrixUnite: 150},
		{Nom: "Thé", Quantite: 1, PrixUnite: 200},
	})
	if err != nil {
		t.Fatalf("erreur inattendue: %v", err)
	}
	if c.ID == "" {
		t.Fatal("attendu un ID généré, obtenu vide")
	}
	if c.Statut != StatutEnCours {
		t.Fatalf("attendu statut en_cours, obtenu %q", c.Statut)
	}
	if c.CreeLe.IsZero() {
		t.Fatal("attendu une date de création, obtenu zéro")
	}
	if c.Montant() != 500 {
		t.Fatalf("attendu montant 500, obtenu %d", c.Montant())
	}
	if repo.created == nil || repo.created.ID != c.ID {
		t.Fatal("attendu la commande persistée via le repository")
	}
}

func TestCreateCommande_ClientRequis(t *testing.T) {
	svc := New(fakeRepo{}, &fakeCommandeRepo{})

	_, err := svc.CreateCommande("", []Produit{{Nom: "Café", Quantite: 1, PrixUnite: 150}})
	if !errors.Is(err, ErrClientRequis) {
		t.Fatalf("attendu ErrClientRequis, obtenu %v", err)
	}
}

func TestCreateCommande_ProduitsVides(t *testing.T) {
	svc := New(fakeRepo{}, &fakeCommandeRepo{})

	_, err := svc.CreateCommande("Alice", nil)
	if !errors.Is(err, ErrProduitsVides) {
		t.Fatalf("attendu ErrProduitsVides, obtenu %v", err)
	}
}

func TestCreateCommande_ProduitInvalide(t *testing.T) {
	svc := New(fakeRepo{}, &fakeCommandeRepo{})

	_, err := svc.CreateCommande("Alice", []Produit{{Nom: "Café", Quantite: 0, PrixUnite: 150}})
	if !errors.Is(err, ErrProduitInvalide) {
		t.Fatalf("attendu ErrProduitInvalide, obtenu %v", err)
	}
}

func TestCreateCommande_ErreurStore(t *testing.T) {
	boom := errors.New("store indisponible")
	svc := New(fakeRepo{}, &fakeCommandeRepo{errCreate: boom})

	_, err := svc.CreateCommande("Alice", []Produit{{Nom: "Café", Quantite: 1, PrixUnite: 150}})
	if !errors.Is(err, boom) {
		t.Fatalf("attendu l'erreur du store, obtenu %v", err)
	}
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
