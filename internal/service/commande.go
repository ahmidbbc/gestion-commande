// Package service holds the domain model and usecases: business orchestration.
package service

import (
	"errors"
	"time"
)

// Domain validation errors for a commande.
var (
	ErrClientRequis    = errors.New("client requis")
	ErrProduitsVides   = errors.New("au moins un produit requis")
	ErrProduitInvalide = errors.New("produit invalide")
)

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

// Produit is a single line item within a commande.
type Produit struct {
	Nom       string
	Quantite  int
	PrixUnite int // in cents, to avoid float rounding issues
}

// Valide reports whether the produit has a name, positive quantity and
// non-negative unit price.
func (p Produit) Valide() bool {
	return p.Nom != "" && p.Quantite > 0 && p.PrixUnite >= 0
}

// Commande is the order domain entity.
type Commande struct {
	ID       string
	Client   string
	Produits []Produit
	Statut   Statut
	CreeLe   time.Time
}

// Montant returns the total amount of the commande in cents.
func (c Commande) Montant() int {
	total := 0
	for _, p := range c.Produits {
		total += p.Quantite * p.PrixUnite
	}
	return total
}

// Valide checks the base invariants of a commande: a client is set and at
// least one valid produit is present.
func (c Commande) Valide() error {
	if c.Client == "" {
		return ErrClientRequis
	}
	if len(c.Produits) == 0 {
		return ErrProduitsVides
	}
	for _, p := range c.Produits {
		if !p.Valide() {
			return ErrProduitInvalide
		}
	}
	return nil
}

// CommandeRepository is the persistence port for commandes. The concrete
// implementation lives in internal/store.
type CommandeRepository interface {
	Save(c Commande) error
	Get(id string) (Commande, error)
	List() ([]Commande, error)
	ListByStatut(statut Statut) ([]Commande, error)
}
