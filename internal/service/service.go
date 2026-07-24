// Package service holds usecases: business orchestration. It depends on
// repository interfaces, never on the api (handler) layer.
package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// ErrStatutInvalide is returned when a requested status is not a known Statut.
var ErrStatutInvalide = errors.New("statut invalide")

// Repository is the data-access port the service depends on. The concrete
// implementation lives in internal/store.
type Repository interface {
	Ping() error
}

// Service orchestrates business logic over its repository dependencies.
type Service struct {
	repo      Repository
	commandes CommandeRepository
}

// New builds a Service with its injected repositories.
func New(repo Repository, commandes CommandeRepository) *Service {
	return &Service{repo: repo, commandes: commandes}
}

// CreateCommande validates the client and produits, stamps the new commande
// with an ID, an "en_cours" status and a creation time, then persists it.
func (s *Service) CreateCommande(client string, produits []Produit) (Commande, error) {
	c := Commande{
		ID:       nouvelID(),
		Client:   client,
		Produits: produits,
		Statut:   StatutEnCours,
		CreeLe:   time.Now(),
	}
	if err := c.Valide(); err != nil {
		return Commande{}, err
	}
	if err := s.commandes.Create(c); err != nil {
		return Commande{}, err
	}
	return c, nil
}

// nouvelID returns a random hex identifier for a new commande.
func nouvelID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

// GetCommande returns the commande matching the given order number (ID).
func (s *Service) GetCommande(id string) (Commande, error) {
	return s.commandes.Get(id)
}

// ListCommandes returns all commandes. When statut is non-empty, only the
// commandes with that status are returned; an unknown status is rejected.
func (s *Service) ListCommandes(statut Statut) ([]Commande, error) {
	if statut == "" {
		return s.commandes.List()
	}
	if !statut.Valide() {
		return nil, ErrStatutInvalide
	}
	return s.commandes.ListByStatut(statut)
}

// UpdateStatut sets a new lifecycle status on an existing commande.
func (s *Service) UpdateStatut(id string, statut Statut) (Commande, error) {
	if !statut.Valide() {
		return Commande{}, ErrStatutInvalide
	}
	c, err := s.commandes.Get(id)
	if err != nil {
		return Commande{}, err
	}
	c.Statut = statut
	if err := s.commandes.Save(c); err != nil {
		return Commande{}, err
	}
	return c, nil
}
