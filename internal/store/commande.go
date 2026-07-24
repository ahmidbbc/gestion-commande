// Package store is a repository: data access only, ZERO business logic.
package store

import (
	"errors"
	"sync"

	"gestion-commande/internal/service"
)

// ErrCommandeIntrouvable is returned when a commande ID has no match.
var ErrCommandeIntrouvable = errors.New("commande introuvable")

// ErrCommandeExiste is returned when creating a commande whose ID already exists.
var ErrCommandeExiste = errors.New("commande déjà existante")

// CommandeStore persists commandes in memory.
type CommandeStore struct {
	mu        sync.RWMutex
	commandes map[string]service.Commande
}

// NewCommandeStore builds an empty in-memory commande repository.
func NewCommandeStore() *CommandeStore {
	return &CommandeStore{commandes: make(map[string]service.Commande)}
}

// Create inserts a new commande, rejecting an ID that already exists.
func (s *CommandeStore) Create(c service.Commande) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.commandes[c.ID]; ok {
		return ErrCommandeExiste
	}
	s.commandes[c.ID] = c
	return nil
}

// Save inserts or replaces a commande.
func (s *CommandeStore) Save(c service.Commande) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.commandes[c.ID] = c
	return nil
}

// Get returns the commande with the given ID or ErrCommandeIntrouvable.
func (s *CommandeStore) Get(id string) (service.Commande, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.commandes[id]
	if !ok {
		return service.Commande{}, ErrCommandeIntrouvable
	}
	return c, nil
}

// List returns all persisted commandes.
func (s *CommandeStore) List() ([]service.Commande, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]service.Commande, 0, len(s.commandes))
	for _, c := range s.commandes {
		out = append(out, c)
	}
	return out, nil
}

// ListByStatut returns the commandes whose status equals statut.
func (s *CommandeStore) ListByStatut(statut service.Statut) ([]service.Commande, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]service.Commande, 0)
	for _, c := range s.commandes {
		if c.Statut == statut {
			out = append(out, c)
		}
	}
	return out, nil
}
