// Package store is a repository: data access only, ZERO business logic.
package store

import "gestion-commande/internal/service"

// Store is the concrete repository implementation.
type Store struct {
	orders []service.Order
}

// New builds a Store seeded with sample orders.
func New() *Store {
	return &Store{
		orders: []service.Order{
			{ID: "CMD-1001", Customer: "Alice Martin", Total: 149.90, Status: "Payée"},
			{ID: "CMD-1002", Customer: "Bruno Da Silva", Total: 79.00, Status: "En cours"},
			{ID: "CMD-1003", Customer: "Chloé Bernard", Total: 320.50, Status: "Expédiée"},
			{ID: "CMD-1004", Customer: "David Nguyen", Total: 24.99, Status: "Annulée"},
		},
	}
}

// Ping reports repository readiness.
func (s *Store) Ping() error {
	return nil
}

// Orders returns the stored customer orders.
func (s *Store) Orders() ([]service.Order, error) {
	return s.orders, nil
}
