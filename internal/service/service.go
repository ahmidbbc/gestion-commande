// Package service holds usecases: business orchestration. It depends on
// repository interfaces, never on the api (handler) layer.
package service

import "errors"

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
