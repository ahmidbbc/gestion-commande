// Package service holds the domain model and usecases: business orchestration.
package service

import "time"

// Statut represents the lifecycle state of a commande.
type Statut string

const (
	StatutEnCours Statut = "en_cours"
	StatutLivree  Statut = "livree"
	StatutAnnulee Statut = "annulee"
)

// Valide reports whether s is a known commande status.
func (s Statut) Valide() bool {
	switch s {
	case StatutEnCours, StatutLivree, StatutAnnulee:
		return true
	default:
		return false
	}
}

// Commande is the order domain entity.
type Commande struct {
	ID     string
	Statut Statut
	CreeLe time.Time
}

// CommandeRepository is the persistence port for commandes. The concrete
// implementation lives in internal/store.
type CommandeRepository interface {
	Save(c Commande) error
	Get(id string) (Commande, error)
	List() ([]Commande, error)
}
